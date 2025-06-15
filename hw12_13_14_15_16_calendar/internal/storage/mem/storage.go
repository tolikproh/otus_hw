package mem

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/pkg/utils"
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
		data: make(map[int]model.Event, 1024),
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
	event.DateTimeStart = event.DateTimeStart.UTC()
	event.DateTimeEnd = event.DateTimeEnd.UTC()

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

	if id > 0 {
		event.ID = id
	} else {
		id = event.ID
	}

	event.DateTimeStart = event.DateTimeStart.UTC()
	event.DateTimeEnd = event.DateTimeEnd.UTC()

	s.data[id] = event

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

	date := utils.GetLastYear(time.Now())

	for _, event := range s.data {
		if event.DateTimeEnd.Compare(date) <= 0 {
			delete(s.data, event.ID)
		}
	}

	return nil
}

func (s *Storage) GetEvents(_ context.Context) ([]model.Event, error) {
	result := make([]model.Event, 0)

	for _, event := range s.data {
		result = append(result, event)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].DateTimeStart.Before(result[j].DateTimeStart)
	})

	return result, nil
}

func (s *Storage) GetEventsByDay(_ context.Context, date time.Time) ([]model.Event, error) {
	start := utils.GetStartOfDay(date)
	end := utils.GetEndOfDay(date)
	return s.findByDateTimeBetween(start, end)
}

func (s *Storage) GetEventsByWeek(_ context.Context, date time.Time) ([]model.Event, error) {
	start := utils.GetStartOfWeek(date)
	end := utils.GetEndOfWeek(date)
	return s.findByDateTimeBetween(start, end)
}

func (s *Storage) GetEventsByMonth(_ context.Context, date time.Time) ([]model.Event, error) {
	start := utils.GetStartOfMonth(date)
	end := utils.GetEndOfMonth(date)
	return s.findByDateTimeBetween(start, end)
}

func (s *Storage) findByDateTimeBetween(start, end time.Time) ([]model.Event, error) {
	result := make([]model.Event, 0)

	for _, v := range s.data {
		if (v.DateTimeEnd.Compare(start) >= 0) && (v.DateTimeEnd.Compare(end) <= 0) {
			result = append(result, v)
		}
	}
	return result, nil
}
