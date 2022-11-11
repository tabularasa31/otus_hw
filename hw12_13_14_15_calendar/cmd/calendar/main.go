package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to init config: %v", err)
		os.Exit(1)
	}

	logg, err := logger.New(conf.Logger.Level)
	if err != nil {
		fmt.Printf("failed to get logger: %v", err)
		os.Exit(1)
	}

	storage := storage.New(conf.Storage.Type)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
