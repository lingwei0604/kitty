//go:build wireinject
// +build wireinject

package wechatcallback

import (
	"git.yingzhongshare.com/mkt/kitty/app/entity"
	"git.yingzhongshare.com/mkt/kitty/app/handlers"
	app "git.yingzhongshare.com/mkt/kitty/app/module"
	"git.yingzhongshare.com/mkt/kitty/app/repository"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	kittyhttp "git.yingzhongshare.com/mkt/kitty/pkg/khttp"
	"git.yingzhongshare.com/mkt/kitty/pkg/ots3"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
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
