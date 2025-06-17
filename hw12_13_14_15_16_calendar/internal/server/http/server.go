package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
)

type Server struct {
	app *app.App
	log *logger.Logger
	cfg *config.Config
	srv *http.Server
}

func New(cfg *config.Config, log *logger.Logger, app *app.App) *Server {
	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.HTTPServer.Host, strconv.Itoa(cfg.HTTPServer.Port)),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return &Server{cfg: cfg, log: log, app: app, srv: srv}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("http server started", "address", s.srv.Addr)

	s.srv.Handler = s.initRoute()
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Debug("http server is shutting down...")
	return s.srv.Shutdown(ctx)
}
