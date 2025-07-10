package http

import (
	"go.uber.org/fx"
)

func ServerModule() fx.Option {
	return fx.Module(
		"http",
		fx.Provide(NewServer),
	)
}
