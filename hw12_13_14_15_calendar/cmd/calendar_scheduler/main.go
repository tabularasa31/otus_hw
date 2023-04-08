package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	proto "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/api"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Notification struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
	UID   string `json:"uid"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/scheduler_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// Configuration
	cfg, err := config.NewSchedulerConfig(configFile)
	failOnError(err, "scheduler config error")

	// AMQP
	mqConn, ch, err := rabbitmq.NewRabbitMQConn(&cfg.AMQPConfig)
	failOnError(err, "...failed to make connection")

	// Close Channel
	defer func() {
		if e := ch.Close(); e != nil {
			failOnError(err, fmt.Sprintf("...failed to close amqp channel, error: %v\n", e))
		}
	}()

	// Close Connection
	defer func() {
		if e := mqConn.Close(); e != nil {
			failOnError(err, fmt.Sprintf("...failed to close amqp connection, error: %v\n", e))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// GRPC client
	conn, err := grpc.Dial("calendar:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	failOnError(err, "failed to connect grpc server")
	defer func() {
		if e := conn.Close(); e != nil {
			failOnError(err, fmt.Sprintf("...failed to close grpc connection, error: %v\n", e))
		}
	}()
	client := proto.NewEventServiceClient(conn)

	// Delete old events Scheduler
	go func() {
		res, e := deleteOldEvents(client)
		if e != nil {
			fmt.Printf("...error while deleting old events: %v\n", e)
		} else {
			fmt.Printf("old events deleted successfully: %s\n", res)
		}
		time.Sleep(24 * time.Hour)
	}()

	// Notifications Scheduler
	for {
		res, e := checkEventsNotifications(client)
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
					cfg.Exchange,   // exchange
					cfg.BindingKey, // routing key
					false,          // mandatory
					false,          // immediate
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

// checkEventsNotifications get events by date notification
func checkEventsNotifications(c proto.EventServiceClient) (*proto.GetEventsResponse, error) {
	t := time.Now().Format("2006-01-02 15:04:05")
	req := &proto.Time{Start: t}

	res, err := c.GetNotificationEvents(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("...error while calling GetNotificationEvents RPC: %w", err)
	}
	return res, nil
}

// deleteOldEvents delete events older 1 year
func deleteOldEvents(c proto.EventServiceClient) (string, error) {
	t := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	req := &proto.Time{Start: t}

	res, err := c.DeleteOldEvents(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("...error while calling DeleteOldEvents RPC: %w", err)
	}
	return res.String(), nil
}

// failOnError wrap error
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
