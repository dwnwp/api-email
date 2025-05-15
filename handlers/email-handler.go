package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dwnwp/api-email/models"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

func ProducerEmail() echo.HandlerFunc {
	return func(c echo.Context) error {

		emailRequest := new(models.MailerRequest)
		if err := c.Bind(emailRequest); err != nil {
			return c.JSON(http.StatusBadRequest, models.RespondError{Error: "error invalid input"})
		}

		// Connect to RabbitMQ
		connString := "amqp://" + os.Getenv("RABBITMQ_USERNAME") + ":" + os.Getenv("RABBITMQ_PASSWORD") + "@localhost:" + os.Getenv("RABBITMQ_PORT") + "/"
		conn, err := amqp.Dial(connString)
		if err != nil {
			return c.JSON(http.StatusOK, models.RespondError{Error: fmt.Sprintf("failed to connect to rabbitmq: %v", err)})
		}
		defer conn.Close()

		// Open channel
		ch, err := conn.Channel()
		if err != nil {
			return c.JSON(http.StatusOK, models.RespondError{Error: fmt.Sprintf("failed to open channel: %v", err)})
		}
		defer ch.Close()

		// Declare queue (will be created if it not exist)
		queueName := "SendEmail"
		_, err = ch.QueueDeclare(
			queueName,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			return c.JSON(http.StatusOK, models.RespondError{Error: fmt.Sprintf("failed to declare queue: %v", err)})
		}

		// Publish message
		body, err := json.Marshal(emailRequest)
		if err != nil {
			log.Fatalf("Failed to marshal message: %v", err)
		}
		err = ch.Publish(
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			return c.JSON(http.StatusOK, models.RespondError{Error: fmt.Sprintf("failed to publish message to rabbitmq: %v", err)})
		}

		log.Println("✅ Message published to queue:", queueName)

		return c.JSON(http.StatusOK, models.RespondMessage{Message: "successfully"})
	}
}
