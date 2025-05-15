package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dwnwp/api-email/models"
	"github.com/dwnwp/api-email/services"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
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
		conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel (Consumer)
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	queueName := "SendEmail"
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
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

	forever := make(chan bool)
	mailer := services.NewMailer()

	go func() {
		for d := range msgs {
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

			if err := mailer.Send(msg.From, msg.To, msg.Subject, bodyTemplate) ; err != nil {
				log.Printf("Failed to send an Email: %v", err)
			}
		}
	}()

	<-forever
}
