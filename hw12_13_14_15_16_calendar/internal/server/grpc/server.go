package grpc

import (
	"context"
	"net"
	"strconv"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/pb"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/server/grpc/interceptor"
	grpclib "google.golang.org/grpc"
)

type Server struct {
	addr    string
	log     *logger.Logger
	cfg     *config.Config
	handler *eventDataHandler
	srv     *grpclib.Server
}

func New(cfg *config.Config, log *logger.Logger, app *app.App) *Server {
	ic := interceptor.New(log)

	addr := net.JoinHostPort(cfg.GRPCServer.Host, strconv.Itoa(cfg.GRPCServer.Port))

	server := &Server{
		addr:    addr,
		log:     log,
		cfg:     cfg,
		handler: newEventDataHandler(log, app),
		srv:     grpclib.NewServer(grpclib.UnaryInterceptor(ic.Logging())),
	}

	pb.RegisterCalendarServer(server.srv, server.handler)

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("grpc server started", "address", s.addr)

	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.srv.Serve(lsn)
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Debug("grpc server is shutting down...")

	s.srv.GracefulStop()
	s.log.Info("grpc server stopped")
	return nil
}
