package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dwnwp/api-email/models"
	"github.com/dwnwp/api-email/services"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

func ProducerEmail(c echo.Context) error {
	emailRequest := new(models.MailerRequest)
	if err := c.Bind(emailRequest); err != nil {
		return c.JSON(http.StatusBadRequest, models.RespondError{Error: "error invalid input"})
	}

	// Connect to RabbitMQ
	connString := "amqp://" + os.Getenv("RABBITMQ_USERNAME") + ":" + os.Getenv("RABBITMQ_PASSWORD") + "@localhost:" + os.Getenv("RABBITMQ_PORT") + "/"
	rabbitmq, err := services.ConnectToRabbitMQ(connString)
	if err != nil {
		return c.JSON(http.StatusOK, models.RespondError{Error: err.Error()})
	}
	defer rabbitmq.DisconnectFromRabbitMQ()

	// Declare queue (will be created if it not exist)
	queueName := "SendEmail"
	if _, err := rabbitmq.Channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return c.JSON(http.StatusOK, models.RespondError{Error: err.Error()})
	}

	// Publish message
	body, err := json.Marshal(emailRequest)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}
	if err = rabbitmq.Channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return c.JSON(http.StatusOK, models.RespondError{Error: fmt.Sprintf("failed to publish message to rabbitmq: %v", err)})
	}

	log.Println("✅ Message published to queue:", queueName)
	return c.JSON(http.StatusOK, models.RespondMessage{Message: "successfully"})
}
