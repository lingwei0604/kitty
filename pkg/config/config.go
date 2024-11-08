package config

import (
	"time"

	corecontract "github.com/DoNewsCode/core/contract"
	"github.com/knadh/koanf"
	"github.com/lingwei0604/kitty/pkg/contract"
)

type Env string
type AppName string

func (a AppName) String() string {
	return string(a)
}

func (e Env) IsLocal() bool {
	return e == "local"
}

func (e Env) IsTesting() bool {
	return e == "testing"
}

func (e Env) IsDev() bool {
	return e == "dev"
}

func (e Env) IsProd() bool {
	return e == "prod"
}

func (e Env) String() string {
	return string(e)
}

func ProvideEnv(conf contract.ConfigReader) Env {
	return Env(conf.String("env"))
}

func ProvideAppName(conf contract.ConfigReader) AppName {
	return AppName(conf.String("name"))
}

var _ corecontract.ConfigAccessor = (*KoanfAdapter)(nil)

type KoanfAdapter struct {
	k *koanf.Koanf
}

func (k *KoanfAdapter) Cut(s string) contract.ConfigReader {
	cut := k.k.Cut("global")
	cut.Merge(k.k.Cut(s))
	return NewKoanfAdapter(cut)
}

func NewKoanfAdapter(k *koanf.Koanf) *KoanfAdapter {
	return &KoanfAdapter{k}
}

func (k *KoanfAdapter) String(s string) string {
	return k.k.String(s)
}

func (k *KoanfAdapter) Int(s string) int {
	return k.k.Int(s)
}

func (k *KoanfAdapter) Strings(s string) []string {
	return k.k.Strings(s)
}

func (k *KoanfAdapter) Bool(s string) bool {
	return k.k.Bool(s)
}

func (k *KoanfAdapter) Get(s string) interface{} {
	return k.k.Get(s)
}

func (k *KoanfAdapter) Float64(s string) float64 {
	return k.k.Float64(s)
}

func (k *KoanfAdapter) Duration(s string) time.Duration {
	return k.k.Duration(s)
}

func (k *KoanfAdapter) Unmarshal(path string, o interface{}) error {
	return k.k.Unmarshal(path, o)
}
