package main

import (
	"flag"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/rabbitmq"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/sender_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// Configuration
	cfg, err := config.NewSenderConfig(configFile)
	failOnError(err, "sender config error")
	_ = cfg

	// AMPQ consumer
	mqConn, ch, err := rabbitmq.NewRabbitMQConn(&cfg.AMQPConfig)

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer mqConn.Close()

	q, err := ch.QueueDeclare(
		cfg.Queue, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "...failed to declare a queue")

	err = ch.QueueBind(
		q.Name,         // queue name
		cfg.BindingKey, // routing key
		cfg.Exchange,   // exchange
		false,
		nil,
	)
	failOnError(err, "...failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,   // queue
		"sender", // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "...failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] Recieved: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
