package internalhttp

import (
	"context"
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/config"
	"net/http"
	"time"
)

type Server struct {
	conf       config.HTTPConfig
	app        Application
	logg       Logger
	httpServer *http.Server
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Application interface {
}

func NewServer(conf config.HTTPConfig, calendar Application, logg Logger) *Server {
	return &Server{conf: conf, app: calendar, logg: logg}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", Hello)
	addr := fmt.Sprintf("%s:%s", s.conf.Host, s.conf.Port)
	http.ListenAndServe(addr, loggingMiddleware(mux))

	<-ctx.Done()
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.Stop(ctx)
		if err != nil {
			s.logg.Error("failed to stop server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func Hello(_ http.ResponseWriter, _ *http.Request) {
	fmt.Println("Hello, world!")
}
