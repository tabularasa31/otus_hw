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
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?

	timeout, addr := handleFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprint(os.Stderr, "...failed to connect, error: %v\n", err)
	} else {
		fmt.Fprint(os.Stderr, "...Connected to %v\n", addr)
	}
	defer func() {
		if e := client.Close(); e != nil {
			fmt.Fprintf(os.Stderr, "...failed to close client, error: %v\n", e)
		}
	}()

	go recieveRoutine(client, cancel)
	go sendRoutine(ctx, client, cancel)

	<-ctx.Done()
}

func handleFlags() (timeout *time.Duration, addr string) {
	timeout = flag.Duration("timeout", 10*time.Second, "Connection timeout in seconds")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalf("there are no any command arguments")
	}

	return timeout, net.JoinHostPort(flag.Arg(0), flag.Arg(1))
}

func recieveRoutine(client TelnetClient, cancel func()) {
	err := client.Receive()
	if err != nil {
		fmt.Fprint(os.Stderr, "recieving error: %v\n", err)
	} else {
		fmt.Fprint(os.Stderr, "...Connection was closed by peer\n")
	}
	defer cancel()
}
func sendRoutine(ctx context.Context, client TelnetClient, cancel func()) {
	err := client.Send()
	if err != nil {
		fmt.Fprint(os.Stderr, "sending error: %v\n", err)
	} else {
		fmt.Fprint(os.Stderr, "...EOF")
	}
	defer cancel()
}
