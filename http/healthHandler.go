package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServiceStatus string

const (
	StatusUP   ServiceStatus = "UP"
	StatusDOWN ServiceStatus = "DOWN"
)

type HealthHandler struct {
	HealthChecks *HealthChecks
}

func NewHealthHandler(hc *HealthChecks) *HealthHandler {
	return &HealthHandler{
		HealthChecks: hc,
	}
}

type ServiceStatusResponse struct {
	Status  ServiceStatus `json:"status"`
	Details interface{}   `json:"details,omitempty"`
	Error   string        `json:"error,omitempty"`
}

func (h *HealthHandler) CustomHealthHandler(c echo.Context) error {
	response := make(map[string]ServiceStatusResponse)

	kafkaStatus, kafkaErr := h.HealthChecks.KafkaCheck.Status()
	if kafkaErr != nil {
		response["kafka"] = ServiceStatusResponse{
			Status:  StatusDOWN,
			Details: kafkaStatus,
			Error:   kafkaErr.Error(),
		}
	} else {
		response["kafka"] = ServiceStatusResponse{
			Status:  StatusUP,
			Details: kafkaStatus,
		}
	}

	postgresStatus, postgresErr := h.HealthChecks.PostgresCheck.Status()
	if postgresErr != nil {
		response["postgres"] = ServiceStatusResponse{
			Status:  StatusDOWN,
			Details: postgresStatus,
			Error:   postgresErr.Error(),
		}
	} else {
		response["postgres"] = ServiceStatusResponse{
			Status:  StatusUP,
			Details: postgresStatus,
		}
	}

	redisStatus, redisErr := h.HealthChecks.RedisCheck.Status()
	if redisErr != nil {
		response["redis"] = ServiceStatusResponse{
			Status:  StatusDOWN,
			Details: redisStatus,
			Error:   redisErr.Error(),
		}
	} else {
		response["redis"] = ServiceStatusResponse{
			Status:  StatusUP,
			Details: redisStatus,
		}
	}

	for _, service := range response {
		if service.Status == StatusDOWN {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	return c.JSON(http.StatusOK, response)
}
