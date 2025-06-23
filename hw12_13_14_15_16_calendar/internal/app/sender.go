package app

import (
	"context"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/brokers/rabbitmq"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
)

type EventSender struct {
	app    *App
	log    *logger.Logger
	broker *rabbitmq.Manager
}

func NewEventSender(app *App, log *logger.Logger, broker *rabbitmq.Manager) *EventSender {
	return &EventSender{app: app, log: log, broker: broker}
}

func (s *EventSender) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.broker.Consume(ctx, cnst.EventQueueName, s.sendNotification)
			if err != nil {
				s.log.Error(err.Error())
			}
		}
	}
}

func (s *EventSender) sendNotification(message []byte) error {
	s.log.Info(string(message))
	return nil
}
