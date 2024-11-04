//go:generate mockery --name=ConfigReader
package contract

import (
	"github.com/DoNewsCode/core/contract"
)

type ConfigReader interface {
	Cut(string) ConfigReader
	contract.ConfigAccessor
}

type Env interface {
	IsLocal() bool
	IsTesting() bool
	IsDev() bool
	IsProd() bool
	String() string
}

type AppName interface {
	String() string
}

type PackageName interface {
	String() string
}
