package module

import (
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/lingwei0604/kitty/app/listener"
	"github.com/lingwei0604/kitty/app/svc"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/event"
	code "github.com/lingwei0604/kitty/pkg/invitecode"
	kittyhttp "github.com/lingwei0604/kitty/pkg/khttp"
	"github.com/lingwei0604/kitty/pkg/kkafka"
	kclient "github.com/lingwei0604/kitty/pkg/kkafka/client"
	logging "github.com/lingwei0604/kitty/pkg/klog"
	"github.com/lingwei0604/kitty/pkg/kmiddleware"
	"github.com/lingwei0604/kitty/pkg/otgorm"
	"github.com/lingwei0604/kitty/pkg/otredis"
	"github.com/lingwei0604/kitty/pkg/ots3"
	"github.com/lingwei0604/kitty/pkg/sms"
	"github.com/lingwei0604/kitty/pkg/wechat"
	kitty "github.com/lingwei0604/kitty/proto"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegermetric "github.com/uber/jaeger-lib/metrics"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	his         metrics.Histogram
	initMetrics sync.Once
)

const maxOpenConns = 20

func provideTokenizer(conf contract.ConfigReader) *code.Tokenizer {
	return code.NewTokenizer(conf.String("salt"))
}

func ProvideHistogramMetrics(appName contract.AppName, env contract.Env) metrics.Histogram {

	return kmiddleware.ProvideHistogramMetrics()
}

func provideKeyManager(appName contract.AppName, env contract.Env) otredis.KeyManager {
	return otredis.NewKeyManager(":", appName.String(), env.String())
}

func ProvideHttpClient(tracer opentracing.Tracer) *kittyhttp.Client {
	return kittyhttp.NewClient(tracer)
}

func providePublisherOptions(tracer opentracing.Tracer, logger log.Logger) []kkafka.PublisherOption {
	return []kkafka.PublisherOption{
		kkafka.PublisherBefore(kkafka.ContextToKafka(tracer, logger)),
	}
}

type producerMiddleware func(operationName string) endpoint.Middleware

func provideProducerMiddleware(tracer opentracing.Tracer, logger log.Logger) producerMiddleware {
	return func(operationName string) endpoint.Middleware {
		return endpoint.Chain(
			kmiddleware.NewAsyncMiddleware(logger),
			kmiddleware.TraceProducer(tracer, operationName),
			kmiddleware.NewTimeoutMiddleware(time.Second),
		)
	}
}

func provideEventBus(factory *kkafka.KafkaFactory, conf contract.ConfigReader, option []kkafka.PublisherOption, mw producerMiddleware) *kclient.EventStore {
	return kclient.NewEventStore(conf.String("kafka.eventBus"), factory, option, mw("kafka.Event"))
}

func provideUserBus(factory *kkafka.KafkaFactory, conf contract.ConfigReader, option []kkafka.PublisherOption, mw producerMiddleware) *kclient.DataStore {
	return kclient.NewDataStore(conf.String("kafka.userBus"), factory, option, mw("kafka.User"))
}

func ProvideKafkaFactory(conf contract.ConfigReader, logger log.Logger) (*kkafka.KafkaFactory, func()) {
	factory := kkafka.NewKafkaFactory(conf.Strings("kafka.brokers"), log.NewNopLogger())
	return factory, func() {
		_ = factory.Close()
	}
}

func ProvideUploadManager(tracer opentracing.Tracer, conf contract.ConfigReader, client contract.HttpDoer) *ots3.Manager {
	return ots3.NewManager(
		conf.String("s3.accessKey"),
		conf.String("s3.accessSecret"),
		conf.String("s3.endpoint"),
		conf.String("s3.region"),
		conf.String("s3.bucket"),
		ots3.WithTracer(tracer),
		ots3.WithHttpClient(client),
		ots3.WithLocationFunc(func(location string) (uri string) {
			u, err := url.Parse(location)
			if err != nil {
				return location
			}
			return fmt.Sprintf(conf.String("s3.cdnUrl"), u.Path[1:])
		}),
	)
}

func ProvideSecurityConfig(conf contract.ConfigReader) *kmiddleware.SecurityConfig {
	return &kmiddleware.SecurityConfig{
		JwtKey: conf.String("security.key"),
		JwtId:  conf.String("security.kid"),
	}
}

func provideWechatConfig(conf contract.ConfigReader, client contract.HttpDoer) *wechat.WechatConfig {
	return &wechat.WechatConfig{
		WechatAccessTokenUrl: conf.String("wechat.wechatAccessTokenUrl"),
		WeChatGetUserInfoUrl: conf.String("wechat.weChatGetUserInfoUrl"),
		AppId:                conf.String("wechat.appId"),
		AppSecret:            conf.String("wechat.appSecret"),
		Client:               client,
	}
}

func provideSmsConfig(doer contract.HttpDoer, conf contract.ConfigReader) *sms.TransportConfig {
	return &sms.TransportConfig{
		Tag:        conf.String("sms.tag"),
		SendUrl:    conf.String("sms.sendUrl"),
		BalanceUrl: conf.String("sms.balanceUrl"),
		UserName:   conf.String("sms.username"),
		Password:   conf.String("sms.password"),
		Client:     doer,
	}
}

func ProvideRedis(logging log.Logger, conf contract.ConfigReader, tracer opentracing.Tracer) (redis.UniversalClient, func()) {
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    conf.Strings("redis.addrs"),
			DB:       conf.Int("redis.database"),
			Password: conf.String("redis.password"),
		})
	client.AddHook(
		otredis.NewHook(tracer, conf.Strings("redis.addrs"),
			conf.Int("redis.database")))
	return client, func() {
		if err := client.Close(); err != nil {
			level.Error(logging).Log("err", err.Error())
		}
	}
}

func ProvideDialector(conf contract.ConfigReader) (gorm.Dialector, error) {
	databaseType := conf.String("gorm.database")
	if databaseType == "mysql" {
		return mysql.Open(conf.String("gorm.dsn")), nil
	}
	if databaseType == "sqlite" {
		return sqlite.Open(conf.String("gorm.dsn")), nil
	}
	return nil, fmt.Errorf("unknow database type %s", databaseType)
}

func ProvideGormConfig(l log.Logger, conf contract.ConfigReader) *gorm.Config {
	return &gorm.Config{
		Logger:                                   &logging.GormLogAdapter{Logging: l},
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.String("name") + "_", // 表名前缀，`User` 的表名应该是 `t_users`
		},
	}
}

func ProvideGormDB(dialector gorm.Dialector, config *gorm.Config, tracer opentracing.Tracer) (*gorm.DB, func(), error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(2)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	otgorm.AddGormCallbacks(db, tracer)
	return db, func() {
		if sqlDb, err := db.DB(); err == nil {
			sqlDb.Close()
		}
	}, nil
}

func ProvideJaegerLogAdapter(l log.Logger) jaeger.Logger {
	return &logging.JaegerLogAdapter{Logging: l}
}

func ProvideDispatcher(ubus listener.UserBus, ebus listener.EventBus) *event.Dispatcher {
	dispatcher := event.Dispatcher{}
	dispatcher.Subscribe(listener.UserChanged{
		Bus: ubus,
	})
	dispatcher.Subscribe(listener.UserCreated{
		Bus: ebus,
	})
	return &dispatcher
}

func ProvideOpentracing(log jaeger.Logger, conf contract.ConfigReader) (opentracing.Tracer, func(), error) {
	cfg := jaegercfg.Configuration{
		ServiceName: conf.String("name"),
		Sampler: &jaegercfg.SamplerConfig{
			Type:  conf.String("jaeger.sampler.type"),
			Param: conf.Float64("jaeger.sampler.param"),
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           conf.Bool("jaeger.log.enable"),
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := log
	jMetricsFactory := jaegermetric.NullFactory

	// Initialize tracer with a logger and a metrics factory
	var (
		canceler io.Closer
		err      error
	)
	tracer, canceler, err := cfg.NewTracer(jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	if err != nil {
		log.Error(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
		return nil, nil, err
	}
	closer := func() {
		if err := canceler.Close(); err != nil {
			log.Error(err.Error())
		}
	}

	return tracer, closer, nil
}

type overallMiddleware func(endpoints svc.Endpoints) svc.Endpoints

func provideModule(
	db *gorm.DB,
	tracer opentracing.Tracer,
	logger log.Logger,
	middleware overallMiddleware,
	server kitty.AppServer,
	appName contract.AppName,
	conf contract.ConfigReader,
	factory *kkafka.KafkaFactory,
) *Module {
	return &Module{
		appName:   appName,
		db:        db,
		logger:    logger,
		tracer:    tracer,
		endpoints: middleware(svc.NewEndpoints(server)),
		conf:      conf,
		factory:   factory,
	}
}
