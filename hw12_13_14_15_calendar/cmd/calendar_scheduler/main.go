package main

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	proto "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/api"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"time"
)

type Notification struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
	UID   string `json:"uid"`
}

func main() {
	// Configuration
	cfg, err := config.NewSchedulerConfig("./config/scheduler_config.yml")
	failOnError(err, "scheduler config error")

	// AMQP
	mqConn, ch, err := rabbitmq.NewRabbitMQConn(cfg.Addr)

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer mqConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// GRPC client
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	failOnError(err, "failed to connect grpc server")
	defer conn.Close()

	client := proto.NewEventServiceClient(conn)

	// Scheduler
	for {
		res, e := checkEvents(client)
		if e != nil {
			fmt.Println(e)
		}
		if res.Events != nil {
			for _, event := range res.Events {
				n := Notification{
					ID:    strconv.Itoa(int(event.Id)),
					Title: event.Title,
					Date:  event.Start,
					UID:   strconv.Itoa(int(event.UserId)),
				}

				send, e := json.Marshal(n)
				if e != nil {
					fmt.Println("...error while marshal json", e)
				}

				err = ch.PublishWithContext(ctx,
					"notifications", // exchange
					"",              // routing key
					false,           // mandatory
					false,           // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        send,
					})
				failOnError(err, "Failed to publish a message")

				log.Printf(" [x] Sent: %v\n", string(send))
			}
		}
		time.Sleep(time.Second)
	}
}

// checkEvents get events by date notification
func checkEvents(c proto.EventServiceClient) (*proto.GetEventsResponse, error) {
	t := time.Now().Format("2006-01-02 15:04:05")
	req := &proto.Time{Start: t}

	res, err := c.GetNotificationEvents(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("...error while calling GetNotificationEvents RPC: %v", err)
	}
	return res, nil
}

// failOnError wrap error
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
