//go:build wireinject
// +build wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"github.com/lingwei0604/kitty/app/entity"
	"github.com/lingwei0604/kitty/app/module"
	"github.com/lingwei0604/kitty/app/repository"
	"github.com/lingwei0604/kitty/pkg/config"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/event"
	"github.com/lingwei0604/kitty/pkg/invitecode"
	kittyhttp "github.com/lingwei0604/kitty/pkg/khttp"
	kclient "github.com/lingwei0604/kitty/pkg/kkafka/client"
	"github.com/lingwei0604/kitty/pkg/ots3"
	"github.com/lingwei0604/kitty/share/handlers"
	"github.com/lingwei0604/kitty/share/internal"
	"github.com/lingwei0604/kitty/share/listener"
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
