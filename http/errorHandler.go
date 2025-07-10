package http

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/qq-mercantil/qq-framework-basic-golang/exceptions"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			errorType := reflect.TypeOf(err).Elem().Name()

			result := map[string]any{
				"error":   errorType,
				"message": err.Error(),
			}

			switch err.(type) {
			case *exceptions.InvalidPropsException, *exceptions.ServiceException, *exceptions.BusinessException:
				return c.JSON(http.StatusBadRequest, result)
			case *exceptions.RepositoryNoDataFoundException:
				return c.JSON(http.StatusNotFound, result)
			default:
				result["error"] = "InternalServerError"
				return c.JSON(http.StatusInternalServerError, result)
			}
		}
		return nil
	}
}
