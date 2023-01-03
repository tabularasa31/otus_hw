package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout, addr := handleFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := NewTelnetClient(ctx, addr, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprint(os.Stderr, "...failed to connect, error: %w\n", err)
	} else {
		fmt.Fprint(os.Stderr, "...Connected to %w\n", addr)
	}
	defer func() {
		if e := client.Close(); e != nil {
			fmt.Fprintf(os.Stderr, "...failed to close client, error: %v\n", e)
		}
	}()

	go receiveRoutine(ctx, client, cancel)
	go sendRoutine(ctx, client, cancel)

	<-ctx.Done()
}

func handleFlags() (timeout *time.Duration, addr string) {
	timeout = flag.Duration("timeout", 10*time.Second, "Connection timeout in seconds")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalf("no any command arguments")
	}

	return timeout, net.JoinHostPort(flag.Arg(0), flag.Arg(1))
}

func receiveRoutine(ctx context.Context, client TelnetClient, cancel func()) {
	select {
	case <-ctx.Done():
		cancel()
	default:
	}
	if err := client.Receive(); err != nil {
		fmt.Fprint(os.Stderr, "receiving error: %w\n", err)
	} else {
		fmt.Fprint(os.Stderr, "...Connection was closed by peer\n")
	}
	defer cancel()
}

func sendRoutine(ctx context.Context, client TelnetClient, cancel func()) {
	select {
	case <-ctx.Done():
		cancel()
	default:
	}
	if err := client.Send(); err != nil {
		fmt.Fprint(os.Stderr, "sending error: %w\n", err)
	} else {
		fmt.Fprint(os.Stderr, "...EOF")
	}
	defer cancel()
}
