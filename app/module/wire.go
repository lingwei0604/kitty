//go:build wireinject
// +build wireinject

package module

import (
	"git.yingzhongshare.com/mkt/kitty/app/entity"
	"git.yingzhongshare.com/mkt/kitty/app/handlers"
	"git.yingzhongshare.com/mkt/kitty/app/listener"
	"git.yingzhongshare.com/mkt/kitty/app/repository"
	"git.yingzhongshare.com/mkt/kitty/pkg/config"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/event"
	kittyhttp "git.yingzhongshare.com/mkt/kitty/pkg/khttp"
	kclient "git.yingzhongshare.com/mkt/kitty/pkg/kkafka/client"
	"git.yingzhongshare.com/mkt/kitty/pkg/otredis"
	"git.yingzhongshare.com/mkt/kitty/pkg/ots3"
	"git.yingzhongshare.com/mkt/kitty/pkg/sms"
	"git.yingzhongshare.com/mkt/kitty/pkg/wechat"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	ProvideDialector,
	ProvideGormConfig,
	ProvideGormDB,
)

var NameAndEnvSet = wire.NewSet(
	config.ProvideAppName,
	config.ProvideEnv,
	wire.Bind(new(contract.Env), new(config.Env)),
	wire.Bind(new(contract.AppName), new(config.AppName)),
)

var OpenTracingSet = wire.NewSet(
	ProvideJaegerLogAdapter,
	ProvideOpentracing,
)

var AppServerSet = wire.NewSet(
	provideSmsConfig,
	DbSet,
	OpenTracingSet,
	NameAndEnvSet,
	provideKeyManager,
	ProvideHttpClient,
	ProvideUploadManager,
	ProvideDispatcher,
	ProvideRedis,
	provideWechatConfig,
	provideUserBus,
	providePublisherOptions,
	ProvideKafkaFactory,
	provideEventBus,
	wechat.NewWechaterFactory,
	wechat.NewWechaterFacade,
	sms.NewTransportFactory,
	sms.NewSenderFacade,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	repository.NewExtraRepo,
	repository.NewUniqueID,
	repository.NewUserCache,
	handlers.NewAppService,
	handlers.ProvideAppServer,
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.Keyer), new(otredis.KeyManager)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(listener.UserBus), new(*kclient.DataStore)),
	wire.Bind(new(listener.EventBus), new(*kclient.EventStore)),
	wire.Bind(new(contract.Dispatcher), new(*event.Dispatcher)),
	wire.Bind(new(wechat.Wechater), new(*wechat.WechaterFacade)),
	wire.Bind(new(contract.SmsSender), new(*sms.SenderFacade)),
	wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(handlers.CodeRepository), new(*repository.CodeRepo)),
	wire.Bind(new(handlers.FileRepository), new(*repository.FileRepo)),
	wire.Bind(new(entity.IDAssigner), new(*repository.UniqueID)),
	wire.Bind(new(handlers.PreAllocator), new(*repository.UniqueID)),
	wire.Bind(new(handlers.UserCache), new(*repository.UserCache)),
)

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	panic(wire.Build(
		AppServerSet,
		ProvideSecurityConfig,
		ProvideHistogramMetrics,
		provideEndpointsMiddleware,
		provideProducerMiddleware,
		provideModule,
	))
}
