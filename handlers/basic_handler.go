package handlers

import (
	"net/http"

	"github.com/dwnwp/api-email/models"
	"github.com/labstack/echo/v4"
)

func Healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, models.RespondMessage{Message: "healthy!"})
}