package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
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
	server := internalhttp.NewServer(cfg, log, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("http server stop", "error", err.Error())
			return
		}
	}()

	if err := server.Start(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("http server stopped")
			return
		}
		log.Error("failed start http server", "error", err)
		return
	}
}
