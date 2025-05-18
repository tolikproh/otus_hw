package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	mem "github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage/memory"
	psql "github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage/postgres"
)

type Storage interface {
	Close() error
	CreateEvent(ctx context.Context, event model.Event) (int, error)
	UpdateEvent(ctx context.Context, id int, event model.Event) error
	DeleteEvent(ctx context.Context, id int) error
	DeleteEventsOldThenLastYear(ctx context.Context) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	GetEventsByLastDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetEventsByLastWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetEventsByLastMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}

func New(ctx context.Context, cfg *config.Config, log *logger.Logger) (Storage, error) {
	switch cfg.Storage.Type {
	case cnst.StorageTypePostgres:
		return psql.New(ctx, cfg, log)

	case cnst.StorageTypeMemory:
		return mem.New(cfg, log), nil

	default:
		return nil, fmt.Errorf(`storage type "%s" not supported`, cfg.Storage.Type)
	}
}
