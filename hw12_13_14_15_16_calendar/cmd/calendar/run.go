package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	internalgrpc "github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/server/grpc"
	internalhttp "github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/server/http"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
)

func run() {
	pathConf := filepath.Dir(configFile)
	fileConf := filepath.Base(configFile)

	cfg := config.NewConfig(pathConf, fileConf)
	log := logger.New(cfg.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := storage.New(ctx, cfg, log)
	if err != nil {
		log.Error("storage init", "error", err.Error())
		return
	}

	calendar := app.New(log, storage)
	serverHTTP := internalhttp.New(cfg, log, calendar)
	serverGRPC := internalgrpc.New(cfg, log, calendar)

	log.Debug("set loger level: " + cfg.Logger.Level)
	log.Debug("set storage type: " + cfg.Storage.Type)
	log.Debug("storage connection: " + cfg.Storage.Conn)

	go func() {
		if err := serverHTTP.Start(ctx); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("http server stopped")
				return
			}
			log.Error("failed start http server", "error", err)
		}
	}()

	go func() {
		if err := serverGRPC.Start(ctx); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("grpc server stopped")
				return
			}
			log.Error("failed start grpc server", "error", err)
		}
	}()

	<-ctx.Done()

	if err := serverHTTP.Stop(ctx); err != nil {
		log.Error("http server stop", "error", err.Error())
	}

	if err := serverGRPC.Stop(ctx); err != nil {
		log.Error("grpc server stop", "error", err.Error())
	}
}
