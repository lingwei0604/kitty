//go:build wireinject
// +build wireinject

package wechatcallback

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"github.com/lingwei0604/kitty/app/entity"
	"github.com/lingwei0604/kitty/app/handlers"
	app "github.com/lingwei0604/kitty/app/module"
	"github.com/lingwei0604/kitty/app/repository"
	"github.com/lingwei0604/kitty/pkg/contract"
	kittyhttp "github.com/lingwei0604/kitty/pkg/khttp"
	"github.com/lingwei0604/kitty/pkg/ots3"
)

var WechatServerSet = wire.NewSet(
	NewHandler,
	app.DbSet,
	app.OpenTracingSet,
	app.NameAndEnvSet,
	app.ProvideHttpClient,
	app.ProvideUploadManager,
	app.ProvideDispatcher,
	app.ProvideRedis,
	app.ProvideKafkaFactory,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	repository.NewExtraRepo,
	repository.NewUniqueID,
	handlers.NewAppService,
	handlers.ProvideAppServer,
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(entity.IDAssigner), new(*repository.UniqueID)),
	wire.Bind(new(handlers.PreAllocator), new(*repository.UniqueID)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Struct(new(Module), "Handler"),
)

func injectModule(reader contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	panic(wire.Build(
		WechatServerSet,
	))
}
