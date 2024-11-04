//go:build wireinject
// +build wireinject

package module

import (
	"git.yingzhongshare.com/mkt/kitty/app/entity"
	"git.yingzhongshare.com/mkt/kitty/app/module"
	"git.yingzhongshare.com/mkt/kitty/app/repository"
	"git.yingzhongshare.com/mkt/kitty/pkg/config"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/event"
	"git.yingzhongshare.com/mkt/kitty/pkg/invitecode"
	kittyhttp "git.yingzhongshare.com/mkt/kitty/pkg/khttp"
	kclient "git.yingzhongshare.com/mkt/kitty/pkg/kkafka/client"
	"git.yingzhongshare.com/mkt/kitty/pkg/ots3"
	"git.yingzhongshare.com/mkt/kitty/share/handlers"
	"git.yingzhongshare.com/mkt/kitty/share/internal"
	"git.yingzhongshare.com/mkt/kitty/share/listener"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
)

var ShareServiceSet = wire.NewSet(
	module.DbSet,
	module.OpenTracingSet,
	module.NameAndEnvSet,
	module.ProvideSecurityConfig,
	module.ProvideKafkaFactory,
	module.ProvideHistogramMetrics,
	module.ProvideHttpClient,
	module.ProvideUploadManager,
	repository.NewUserRepo,
	repository.NewRelationRepo,
	repository.NewFileRepo,
	repository.NewUniqueID,
	ProvideRedis,
	provideTokenizer,
	providePublisherOptions,
	provideInvitationCodeBus,
	provideDispatcher,
	internal.NewXTaskRequester,
	handlers.NewShareService,
	handlers.ProvideShareServer,
	wire.Struct(new(internal.InvitationManagerFactory), "*"),
	wire.Struct(new(internal.InvitationManagerFacade), "*"),
	wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(internal.RelationRepository), new(*repository.RelationRepo)),
	wire.Bind(new(handlers.InvitationManager), new(*internal.InvitationManagerFacade)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(internal.EncodeDecoder), new(*invitecode.Tokenizer)),
	wire.Bind(new(contract.Dispatcher), new(*event.Dispatcher)),
	wire.Bind(new(listener.InvitationCodeBus), new(*kclient.DataStore)),
	wire.Bind(new(entity.IDAssigner), new(*repository.UniqueID)),
)

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	panic(wire.Build(
		ShareServiceSet,
		provideEndpointsMiddleware,
		provideEndpoints,
		provideHttp,
		provideGrpc,
		provideKafkaServer,
		provideProducerMiddleware,
		provideModule,
	))
}
