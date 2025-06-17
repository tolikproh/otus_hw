package mem

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
)

func TestStorage(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		storage := New(nil, nil)

		ctx := context.Background()
		emptyEventsList := make([]model.Event, 0)
		events, err := storage.GetEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, emptyEventsList, events)
	})

	t.Run("use simple script", func(t *testing.T) {
		storage := New(nil, nil)

		ctx := context.Background()
		now := time.Now()
		duration := 24 * time.Hour

		event := model.Event{
			Title:                "TestEventTitle",
			DateTimeStart:        now.AddDate(0, 0, -1),
			DateTimeEnd:          now,
			Description:          "TestEventDescription",
			UserID:               1,
			NotificationDuration: &duration,
		}

		id, err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)

		events, err := storage.GetEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(events), 1)
		require.Equal(t, events[0].Title, event.Title)

		event.Title = "NewTestEventTitle"
		err = storage.UpdateEvent(ctx, id, event)
		require.NoError(t, err)

		events, err = storage.GetEvents(ctx)
		require.NoError(t, err)
		require.Equal(t, len(events), 1)
		require.Equal(t, events[0].Title, event.Title)

		err = storage.DeleteEvent(ctx, id)
		require.NoError(t, err)

		err = storage.UpdateEvent(ctx, id, event)
		require.Equal(t, ErrEventNotFoundError, err)
	})

	t.Run("multithreading", func(t *testing.T) {
		storage := New(nil, nil)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		ctx := context.Background()
		event := model.Event{
			Title:       "TestEventTitle",
			Description: "TestEventDescription",
		}

		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				storage.CreateEvent(ctx, event)
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				storage.CreateEvent(ctx, event)
			}
		}()

		wg.Wait()
	})
}
