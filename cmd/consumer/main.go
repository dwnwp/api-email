package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwnwp/api-email/models"
	"github.com/dwnwp/api-email/services"
	"github.com/joho/godotenv"
)

func init() {
	wd, _ := os.Getwd()
	if err := godotenv.Load(wd + "/.env"); err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func main() {
	// Connect to RabbitMQ
	connString := "amqp://" + os.Getenv("RABBITMQ_USERNAME") + ":" + os.Getenv("RABBITMQ_PASSWORD") + "@localhost:" + os.Getenv("RABBITMQ_PORT") + "/"
	rabbitmq, err := services.ConnectToRabbitMQ(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.DisconnectFromRabbitMQ()

	// Declare the queue
	queueName := "SendEmail"
	if _, err := rabbitmq.Channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages
	messages, err := rabbitmq.Channel.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("🟢 Waiting for messages...")

	mailer := services.NewMailer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range messages {
			var msg models.MailerRequest
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("❌ Error decoding JSON: %v", err)
				continue
			}

			bodyTemplate := models.CreateMailBodyTemplate(msg.BodySubject, msg.BodyContent)

			fmt.Println("📩 Received Email:")
			fmt.Println("From:   ", msg.From)
			fmt.Println("To:     ", msg.To)
			fmt.Println("Subject:", msg.Subject)
			fmt.Println("BodySubject:   ", msg.BodySubject)
			fmt.Println("BodyContent:   ", msg.BodyContent)
			fmt.Println("--------")

			if err := mailer.Send(msg.From, msg.To, msg.Subject, bodyTemplate); err != nil {
				log.Printf("Failed to send an Email: %v", err)
			}
		}
	}()

	<-sigChan
	fmt.Println("Gracefully stop consuming.")
}
