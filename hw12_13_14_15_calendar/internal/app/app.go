package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	v1 "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/http/v1"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo/memoryrepo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo/postgresrepo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/httpserver"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/storage/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {

	// Logger
	logg := logger.New(cfg.Logger.Level)

	// EventRepo
	r := repo(cfg)

	// Use case
	eventUseCase := usecase.New(r)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, logg, *eventUseCase)
	httpServer := httpserver.New(handler, cfg.HTTP, logg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logg.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		logg.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		logg.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func repo(cfg *config.Config) usecase.EventRepo {
	if cfg.Storage.Type == "postgres" {
		pg, err := postgres.New(cfg.Storage.Dsn)
		if err != nil {
			log.Fatal(fmt.Errorf("app - Run - repo - postgres.New: %w", err))
		}
		defer pg.Close()
		return postgresrepo.New(pg)
	}
	return memoryrepo.New()
}
