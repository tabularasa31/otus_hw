package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	v1 "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/http/v1"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo/memoryrepo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo/postgresrepo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/grpcserver"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/httpserver"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/storage/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// Logger
	logg := logger.New(cfg.Logger.Level)

	// EventRepo
	var repo usecase.EventRepo
	if cfg.Storage.Type == "postgres" {
		pg, err := postgres.New(cfg)
		if err != nil {
			log.Fatal(fmt.Errorf("app - Run - repo - postgres.New: %w", err))
		}
		defer pg.Close()
		repo = postgresrepo.New(pg)
	} else {
		repo = memoryrepo.New()
	}

	// Use case
	eventUseCase := usecase.New(repo)

	// RabbitMQ RPC Server
	//rmqRouter := amqprpc.NewRouter(eventUseCase)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, logg, *eventUseCase)
	httpServer := httpserver.New(handler, cfg.HTTP, logg)

	// GRPC Server
	lis, err := net.Listen("tcp", cfg.GRPC.Addr)
	if err != nil {
		logg.Fatal(fmt.Errorf("app - Run - net.Listen: %w", err))
	}
	defer func() {
		if e := lis.Close(); e != nil {
			logg.Fatal("...failed to close client, error: %v\n", e)
		}
	}()

	grpcServer := grpcserver.New(lis, logg, eventUseCase)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logg.Info("app - Run - signal: " + s.String())
	case e := <-httpServer.Notify():
		logg.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", e))
	}

	// Shutdown
	errServer := httpServer.Shutdown()
	if errServer != nil {
		logg.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	grpcServer.Shutdown()
}
