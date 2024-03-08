//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/peifengll/go_809_converter/converter/handlers/converters"
)

var repoSet = wire.NewSet()

var converterSet = wire.Build(
	converters.NewRequestConverters,
)

func NewWire() {
	panic(wire.Build())
}
