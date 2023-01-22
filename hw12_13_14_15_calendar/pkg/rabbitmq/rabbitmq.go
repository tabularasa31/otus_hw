package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// Initialize new RabbitMQ connection
func NewRabbitMQConn(addr string) (*amqp.Connection, *amqp.Channel, error) {
	mqConn, err := amqp.Dial(addr)
	failOnError(err, "...failed to connect to RabbitMQ")

	ch, err := mqConn.Channel()
	failOnError(err, "...failed to open a channel")

	err = ch.ExchangeDeclare(
		"events", // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
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
