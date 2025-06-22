package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/brokers/rabbitmq"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
)

type EventScheduler struct {
	app    *App
	log    *logger.Logger
	broker *rabbitmq.Manager
}

func NewEventScheduler(app *App, log *logger.Logger, broker *rabbitmq.Manager) *EventScheduler {
	return &EventScheduler{app: app, log: log, broker: broker}
}

func (s *EventScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.sentEvents(ctx)
			if err != nil {
				s.log.Error(err.Error())
			}

			err = s.deleteOldEvents(ctx)
			if err != nil {
				s.log.Error(err.Error())
			}
		}
	}
}

func (s *EventScheduler) sentEvents(ctx context.Context) error {
	events, err := s.app.GetEventsByDay(ctx)
	if err != nil {
		return err
	}

	for _, event := range events {
		b, err := json.Marshal(&event)
		if err != nil {
			return err
		}

		err = s.broker.Send(b, cnst.EventQueueName, "application/json")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EventScheduler) deleteOldEvents(ctx context.Context) error {
	err := s.app.DeleteEventsOld(ctx)
	if err != nil {
		return err
	}

	return nil
}
