package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/brokers/rabbitmq"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
)

func run() {
	pathConf := filepath.Dir(configFile)
	fileConf := filepath.Base(configFile)

	cfg := config.NewConfig(pathConf, fileConf)
	log := logger.New(cfg.Logger.Level)

	ctx, cancel := context.WithCancel(context.Background())
	go watchExitSignals(cancel)

	storage, err := storage.New(ctx, cfg, log)
	if err != nil {
		log.Error("storage init", "error", err.Error())
		return
	}

	broker, err := rabbitmq.New(cfg, log)
	if err != nil {
		log.Error("brocker rabbitmq init", "error", err.Error())
		return
	}

	calendar := app.New(log, storage)
	eventScheduler := app.NewEventScheduler(calendar, log, broker)
	go eventScheduler.Start(ctx)

	log.Info("calendar_scheduler service is running...")

	<-ctx.Done()

	log.Info("calendar_scheduler service was stopped")
}

func watchExitSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	cancel()
}
