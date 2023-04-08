package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
)

// Initialize new RabbitMQ connection
func NewRabbitMQConn(cfg *config.AMQPConfig) (*amqp.Connection, *amqp.Channel, error) {
	mqConn, err := amqp.Dial(cfg.Addr)
	failOnError(err, "...failed to connect to RabbitMQ")

	ch, err := mqConn.Channel()
	failOnError(err, "...failed to open a channel")

	err = ch.ExchangeDeclare(
		cfg.Exchange,     // name
		cfg.ExchangeType, // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "...failed to declare an exchange")

	return mqConn, ch, nil
}

// failOnError wrap error
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
