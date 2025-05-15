package handlers

import (
	"net/http"

	"github.com/dwnwp/api-email/models"
	"github.com/labstack/echo/v4"
)

func ProducerEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO
		return c.JSON(http.StatusOK, models.RespondMessage{Message: "successfully"})
	}
}