package handlers

import (
	"net/http"

	"github.com/dwnwp/api-email/models"
	"github.com/dwnwp/api-email/services"
	"github.com/labstack/echo/v4"
)

func ProducerEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		
		emailRequest := new(models.MailerRequest)
		if err := c.Bind(emailRequest); err != nil {
			return c.JSON(http.StatusBadRequest, models.RespondError{Error: "error invalid input"})
		}

		mailer := services.NewMailer()
		if err := mailer.Send(emailRequest.From, emailRequest.To, emailRequest.Subject, emailRequest.Body) ; err != nil {
			return c.JSON(http.StatusBadRequest, models.RespondError{Error: err.Error()})
		}
		
		return c.JSON(http.StatusOK, models.RespondMessage{Message: "successfully"})
	}
}