//go:build wireinject
// +build wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/preload/handlers"
)

func injectModule(conf contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	panic(wire.Build(handlers.NewService, provideModule))
}
