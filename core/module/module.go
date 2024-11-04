package module

import (
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"github.com/go-kit/kit/log"
)

type Module struct {
	Conf   contract.ConfigReader
	Logger log.Logger
}

func New(cfgFile string) *Module {
	conf, err := ProvideConfig(cfgFile)
	if err != nil {
		panic(err)
	}
	logger := ProvideLogger(conf)
	return &Module{
		Conf:   conf,
		Logger: logger,
	}
}
