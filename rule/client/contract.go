package client

import (
	kconf "git.yingzhongshare.com/mkt/kitty/pkg/config"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/rule/dto"
)

type Tenanter interface {
	Payload(pl *dto.Payload) (contract.ConfigReader, error)
	Tenant(tenant *kconf.Tenant) (contract.ConfigReader, error)
}

type Engine interface {
	Of(ruleName string) Tenanter
}
