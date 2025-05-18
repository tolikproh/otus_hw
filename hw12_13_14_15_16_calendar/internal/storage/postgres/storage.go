package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
)

type Storage struct {
	cfg *config.Config
	log *logger.Logger
	db  *sql.DB
}

func New(ctx context.Context, cfg *config.Config, log *logger.Logger) (*Storage, error) {
	db, err := sql.Open("pgx", cfg.Storage.Conn)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Info("service connected to database")

	return &Storage{cfg, log, db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) (int, error) {
	args := []interface{}{
		event.Title,
		event.DateTimeStart,
		event.DateTimeEnd,
		event.Description,
		event.UserID,
		event.NotificationDuration,
	}
	query := `INSERT INTO event (title, date_time_start, date_time_end, description, user_id, notification_duration)
			  VALUES($1, $2, $3, $4, $5, $6)
			  RETURNING id`

	var id int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id int, event model.Event) error {
	args := []interface{}{
		event.Title,
		event.DateTimeStart,
		event.DateTimeEnd,
		event.Description,
		event.UserID,
		event.NotificationDuration,
		id,
	}
	query := `UPDATE event
			  SET title = $1,
				  date_time_start = $2,
				  date_time_end = $3,
				  description = $4,
				  user_id = $5,
				  notification_duration = $6
			  WHERE id = $7;`

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	query := `DELETE FROM event	WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEventsOldThenLastYear(ctx context.Context) error {
	query := `DELETE FROM event	WHERE date_time_end < date_trunc('minute', now() - interval '1 year')`
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]model.Event, error) {
	query := `SELECT id, title, date_time_start, date_time_end, description, user_id, notification_duration
			  FROM event
		      ORDER BY date_time_start`
	return s.prepareEvents(ctx, query)
}

func (s *Storage) GetEventsByLastDay(ctx context.Context, date time.Time) ([]model.Event, error) {
	query := `SELECT id, title, date_time_start, date_time_end, description, user_id, notification_duration
			  FROM event
			  WHERE date_time_end > date_trunc('minute', now() - interval '1 day')
		      ORDER BY date_time_start`
	return s.prepareEvents(ctx, query)
}

func (s *Storage) GetEventsByLastWeek(ctx context.Context, date time.Time) ([]model.Event, error) {
	query := `SELECT id, title, date_time_start, date_time_end, description, user_id, notification_duration
			  FROM event
			  WHERE date_time_end > date_trunc('minute', now() - interval '1 week')
		      ORDER BY date_time_start`
	return s.prepareEvents(ctx, query)
}

func (s *Storage) GetEventsByLastMonth(ctx context.Context, date time.Time) ([]model.Event, error) {
	query := `SELECT id, title, date_time_start, date_time_end, description, user_id, notification_duration
			  FROM event
			  WHERE date_time_end > date_trunc('minute', now() - interval '1 month')
		      ORDER BY date_time_start`
	return s.prepareEvents(ctx, query)
}

func (s *Storage) prepareEvents(ctx context.Context, query string) ([]model.Event, error) {
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.Event, 0)
	for rows.Next() {
		var event model.Event
		var notification sql.NullInt64
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.DateTimeStart,
			&event.DateTimeEnd,
			&event.Description,
			&event.UserID,
			&notification,
		)
		if err != nil {
			return nil, err
		}

		if notification.Valid {
			event.NotificationDuration = (*time.Duration)(&notification.Int64)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
