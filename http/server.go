package http

import (
	"context"
	"fmt"
	netHTTP "net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/qq-mercantil/qq-framework-log-golang/logger"
	"go.uber.org/fx"
)

type HttpServer struct {
	AppServer *echo.Echo
	AppGroup  *echo.Group
}

type ServerParams struct {
    fx.In
    Config IHttpProvider
    HealthChecks *HealthChecks `optional:"true"`
    Lifecycle fx.Lifecycle
}

func NewServer(
	params ServerParams,
) *HttpServer {

	log := logger.Get()

	server := echo.New()
	server.Use(logger.EchoLogger)

	appGroup := server.Group(fmt.Sprintf("/%s", os.Getenv("APP_NAME")))

	if params.HealthChecks != nil {
        healthHandler := NewHealthHandler(params.HealthChecks)
        appGroup.GET("/health", healthHandler.CustomHealthHandler)
    }

	appGroup.GET("/alive", func(c echo.Context) error {
		response := map[string]string{
			"message": "Hello, World!",
		}
		return c.JSON(netHTTP.StatusOK, response)
	})

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			port := params.Config.GetHttpPort()

			srv := netHTTP.Server{
				Addr: fmt.Sprintf(":%d", port),
			}
			log.Debugf("Iniciando servidor na porta %d", port)
			go func() {
				if err := server.Start(srv.Addr); err != nil && err != netHTTP.ErrServerClosed {
					log.Errorf("Erro ao iniciar o servidor na porta %d: %s", port, err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Debugf("Parando o servidor...")
			return server.Shutdown(ctx)
		},
	})

	return &HttpServer{
		AppServer: server,
		AppGroup:  appGroup,
	}
}
