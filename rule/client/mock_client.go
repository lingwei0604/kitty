package client

import (
	kconf "git.yingzhongshare.com/mkt/kitty/pkg/config"
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/rule/dto"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
)

// MockRule 是 ofRule 的测试替身。
type MockRule struct {
	f func(pl *dto.Payload) dto.Data
}

func NewMockRule(f func(pl *dto.Payload) dto.Data) *MockRule {
	return &MockRule{f: f}
}

func (m *MockRule) Tenant(tenant *kconf.Tenant) (contract.ConfigReader, error) {
	var payload = dto.FromTenant(tenant)
	return m.Payload(payload)
}

func (m *MockRule) Payload(pl *dto.Payload) (contract.ConfigReader, error) {
	data := m.f(pl)
	c := koanf.New(".")
	err := c.Load(confmap.Provider(data, "."), nil)
	if err != nil {
		return nil, err
	}
	adapter := kconf.NewKoanfAdapter(c)
	return adapter, nil
}

// MockEngine 是 Engine 的测试替身
type MockEngine struct {
	mapping map[string]Tenanter
}

func NewMockEngine(mapping map[string]Tenanter) *MockEngine {
	return &MockEngine{
		mapping: mapping,
	}
}

func (m *MockEngine) Of(ruleName string) Tenanter {
	return m.mapping[ruleName]
}