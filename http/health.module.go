package http

import (
	"go.uber.org/fx"
)

func HealthModule() fx.Option {
	return fx.Module(
		"health",
		fx.Provide(RegisterHealthChecks),
	)
}
