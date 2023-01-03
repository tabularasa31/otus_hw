package httpserver

import (
	"context"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/config"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	logger          logger.Logger
}

func New(handler http.Handler, conf config.HTTP, logg *logger.Logger) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         conf.Addr,
	}
	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          *logg,
	}

	s.Start()

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- http.ListenAndServe(s.server.Addr, loggingMiddleware(s.server.Handler, s.logger))
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
