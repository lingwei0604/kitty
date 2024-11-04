//go:build wireinject
// +build wireinject

package module

import (
	"github.com/DoNewsCode/core/clihttp"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"github.com/lingwei0604/kitty/app/module"
	"github.com/lingwei0604/kitty/pkg/config"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/rule/service"
)

var serviceSet = wire.NewSet(
	provideEtcdClient,
	provideRepository,
	service.ProvideService,
)

func injectModule(reader contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	panic(wire.Build(
		serviceSet,
		newEndpoints,
		module.OpenTracingSet,
		provideModule,
		provideDmpServer,
		provideHistogramMetrics,
		ProvideHTTPClient,
		module.ProvideRedis,
		config.ProvideAppName,
		config.ProvideEnv,
		wire.Bind(new(contract.Env), new(config.Env)),
		wire.Bind(new(contract.AppName), new(config.AppName)),
		wire.Bind(new(contract.HttpDoer), new(*clihttp.Client)),
	))
}
