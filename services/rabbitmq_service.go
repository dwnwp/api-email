package services

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func ConnectToRabbitMQ(connString string) (RabbitMQ, error) {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return RabbitMQ{}, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return RabbitMQ{}, fmt.Errorf("failed to open channel: %v", err)
	}
	return RabbitMQ{Connection: conn, Channel: ch}, nil
}

func (r *RabbitMQ) DisconnectFromRabbitMQ() {
	r.Connection.Close()
	r.Channel.Close()
	log.Printf("RabbitMQ closed connection successfully")
}
