package memorystorage

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
)

var ErrEventNotFoundError = errors.New("event not found")

type Storage struct {
	data   map[int]model.Event
	lastID int
	cfg    *config.Config
	log    *logger.Logger
	mu     sync.RWMutex
}

func New(cfg *config.Config, log *logger.Logger) *Storage {
	return &Storage{
		data: make(map[int]model.Event),
		cfg:  cfg,
		log:  log,
	}
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateEvent(_ context.Context, event model.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++

	event.ID = s.lastID
	s.data[s.lastID] = event

	return s.lastID, nil
}

func (s *Storage) UpdateEvent(_ context.Context, id int, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[id]
	if !ok {
		return ErrEventNotFoundError
	}

	e := event
	s.data[id] = e

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[id]
	if !ok {
		return ErrEventNotFoundError
	}

	delete(s.data, id)

	return nil
}

func (s *Storage) DeleteEventsOldThenLastYear(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for _, event := range s.data {
		if event.DateTimeEnd.Before(now.AddDate(-1, 0, 0)) {
			delete(s.data, event.ID)
		}
	}

	return nil
}

func (s *Storage) GetEvents(_ context.Context) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Event, 0)

	for _, event := range s.data {
		result = append(result, event)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DateTimeStart.Before(result[j].DateTimeStart)
	})

	return result, nil
}

func (s *Storage) GetEventsByLastDay(_ context.Context, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Event, 0)

	year, month, day := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, eventDay := event.DateTimeStart.Date()
		if eventYear == year && eventMonth == month && eventDay == day {
			result = append(result, event)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DateTimeStart.Before(result[j].DateTimeStart)
	})

	return result, nil
}

func (s *Storage) GetEventsByLastWeek(_ context.Context, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Event, 0)

	year, week := date.ISOWeek()
	for _, event := range s.data {
		eventYear, eventWeek := event.DateTimeStart.ISOWeek()
		if eventYear == year && eventWeek == week {
			result = append(result, event)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DateTimeStart.Before(result[j].DateTimeStart)
	})

	return result, nil
}

func (s *Storage) GetEventsByLastMonth(_ context.Context, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Event, 0)

	year, month, _ := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, _ := event.DateTimeStart.Date()
		if eventYear == year && eventMonth == month {
			result = append(result, event)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DateTimeStart.Before(result[j].DateTimeStart)
	})

	return result, nil
}
