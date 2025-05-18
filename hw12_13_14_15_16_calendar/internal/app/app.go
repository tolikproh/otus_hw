package app

import (
	"context"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
)

type App struct {
	log     *logger.Logger
	storage Storage
}

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

func New(log *logger.Logger, storage Storage) *App {
	return &App{log, storage}
}

func (a *App) Close() error {
	return a.storage.Close()
}

func (a *App) CreateEvent(ctx context.Context, id int, title string) (int, error) {
	return a.storage.CreateEvent(ctx, model.Event{ID: id, Title: title})
}

func (a *App) UpdateEvent(ctx context.Context, id int, event model.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) DeleteEventsOldThenLastYear(ctx context.Context) error {
	return a.storage.DeleteEventsOldThenLastYear(ctx)
}

func (a *App) GetEvents(ctx context.Context) ([]model.Event, error) {
	return a.storage.GetEvents(ctx)
}

func (a *App) GetEventsByLastDay(ctx context.Context, date time.Time) ([]model.Event, error) {
	return a.storage.GetEventsByLastDay(ctx, date)
}

func (a *App) GetEventsByLastWeek(ctx context.Context, date time.Time) ([]model.Event, error) {
	return a.storage.GetEventsByLastWeek(ctx, date)
}

func (a *App) GetEventsByLastMonth(ctx context.Context, date time.Time) ([]model.Event, error) {
	return a.storage.GetEventsByLastMonth(ctx, date)
}
