package grpcserver

import (
	"net"
	"time"

	proto "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/api"
	grpcv1 "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/grpc/v1"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	_defaultMaxConnectionIdle = 5 * time.Minute
	_defaultMaxConnectionAge  = 5 * time.Minute
	_defaultTimeout           = 15 * time.Second
	_defaultTime              = 5 * time.Minute
)

type Server struct {
	server *grpc.Server
	notify chan error
	logger logger.Logger
}

func New(lis net.Listener, logg *logger.Logger, u *usecase.EventUseCase) *Server {
	im := NewInterceptorManager(*logg)

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: _defaultMaxConnectionIdle,
		Timeout:           _defaultTimeout,
		MaxConnectionAge:  _defaultMaxConnectionAge,
		Time:              _defaultTime,
	}),
		grpc.UnaryInterceptor(im.Logger),
	)

	reflection.Register(grpcServer)

	s := &Server{
		server: grpcServer,
		notify: make(chan error, 1),
		logger: *logg,
	}

	grpcService := grpcv1.NewCalendarGRPCService(*u, logg)

	proto.RegisterEventServiceServer(grpcServer, grpcService)

	s.Start(lis)

	return s
}

func (s *Server) Start(lis net.Listener) {
	go func() {
		s.notify <- s.server.Serve(lis)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() {
	s.server.GracefulStop()
	s.logger.Info("grpc Server Exited Properly")
}
