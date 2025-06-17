package app

import (
	"context"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
)

type App struct {
	log     *logger.Logger
	storage storage.Storage
}

func New(log *logger.Logger, storage storage.Storage) *App {
	return &App{log, storage}
}

func (a *App) Close() error {
	return a.storage.Close()
}

func (a *App) CreateEvent(ctx context.Context, event model.Event) (int, error) {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, id int, event model.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) DeleteEventsOld(ctx context.Context) error {
	return a.storage.DeleteEventsOldThenLastYear(ctx)
}

func (a *App) GetEvents(ctx context.Context) ([]model.Event, error) {
	return a.storage.GetEvents(ctx)
}

func (a *App) GetEventsByDay(ctx context.Context) ([]model.Event, error) {
	date := time.Now()
	return a.storage.GetEventsByDay(ctx, date)
}

func (a *App) GetEventsByWeek(ctx context.Context) ([]model.Event, error) {
	date := time.Now()
	return a.storage.GetEventsByWeek(ctx, date)
}

func (a *App) GetEventsByMonth(ctx context.Context) ([]model.Event, error) {
	date := time.Now()
	return a.storage.GetEventsByMonth(ctx, date)
}
