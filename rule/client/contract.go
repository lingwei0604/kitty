package client

import (
	kconf "github.com/lingwei0604/kitty/pkg/config"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/rule/dto"
)

type Tenanter interface {
	Payload(pl *dto.Payload) (contract.ConfigReader, error)
	Tenant(tenant *kconf.Tenant) (contract.ConfigReader, error)
}

type Engine interface {
	Of(ruleName string) Tenanter
}
