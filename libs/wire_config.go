// +build wireinject

package libs

import "github.com/google/wire"

func NewConfig() interface{}{
	panic(wire.Build(ProvideConfigType,ProvideConfigFile,ProvideEnvironment,LoadConfig))
	return nil
}
