package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
)

type Manager struct {
	log  *logger.Logger
	conn *amqp.Connection
}

func New(cfg *config.Config, log *logger.Logger) (*Manager, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.Address)
	if err != nil {
		return nil, err
	}

	return &Manager{
		conn: conn,
		log:  log,
	}, nil
}

func (c *Manager) Send(message []byte, queueName, contentType string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        message,
		})
	if err != nil {
		return err
	}

	c.log.Info("Send message in rabbitMQ: " + string(message))

	return nil
}

func (c *Manager) Consume(ctx context.Context, queueName string, do func(message []byte) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			err := do(d.Body)
			if err != nil {
				c.log.Error(err.Error())
			}
		}
	}()

	<-ctx.Done()

	return nil
}

func (c *Manager) Disconnect() error {
	return c.conn.Close()
}
