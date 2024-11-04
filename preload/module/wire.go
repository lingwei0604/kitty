//go:build wireinject
// +build wireinject

package module

import (
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/preload/handlers"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
)

func injectModule(conf contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	panic(wire.Build(handlers.NewService, provideModule))
}
