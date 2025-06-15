package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
)

func testData(now time.Time) []model.Event {
	duration := 30 * time.Minute

	testEvents := []model.Event{
		{
			Title:                "TestEventLastTwoYearOld",
			DateTimeStart:        now.AddDate(-2, 0, 0),
			DateTimeEnd:          now.AddDate(-2, 0, 0).Add(duration),
			Description:          "TestEventLastTwoYearOldDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventLastYearOld",
			DateTimeStart:        now.AddDate(-1, 0, -1),
			DateTimeEnd:          now.AddDate(-1, 0, -1).Add(duration),
			Description:          "TestEventLastYearOldDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventLastMonth",
			DateTimeStart:        now.AddDate(0, -1, 0),
			DateTimeEnd:          now.AddDate(0, -1, 0).Add(duration),
			Description:          "TestEventLastMonthDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventLastWeek",
			DateTimeStart:        now.AddDate(0, 0, -8),
			DateTimeEnd:          now.AddDate(0, 0, -8).Add(duration),
			Description:          "TestEventLastWeekDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventLastDay",
			DateTimeStart:        now.AddDate(0, 0, -1),
			DateTimeEnd:          now.AddDate(0, 0, -1).Add(duration),
			Description:          "TestEventLastDayDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventToday",
			DateTimeStart:        now,
			DateTimeEnd:          now.Add(duration),
			Description:          "TestEventTodayDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventNextDay",
			DateTimeStart:        now.AddDate(0, 0, 1),
			DateTimeEnd:          now.AddDate(0, 0, 1).Add(duration),
			Description:          "TestEventNextDayDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventNextWeek",
			DateTimeStart:        now.AddDate(0, 0, 7),
			DateTimeEnd:          now.AddDate(0, 0, 7).Add(duration),
			Description:          "TestEventNextWeekDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
		{
			Title:                "TestEventNextMonth",
			DateTimeStart:        now.AddDate(0, 1, 0),
			DateTimeEnd:          now.AddDate(0, 1, 0).Add(duration),
			Description:          "TestEventNextMonthDescription",
			UserID:               1,
			NotificationDuration: &duration,
		},
	}

	return testEvents
}

type testHTTP struct {
	ctx            context.Context
	store          storage.Storage
	client         http.Client
	srv            *httptest.Server
	now            time.Time
	eventsExpected []model.Event
	updateID       int
	eventsAll      []model.Event
	expectedTitle  string
}

func testInit() *testHTTP {
	cfg := new(config.Config)
	cfg.Logger.Level = cnst.LoggerLevelDebug
	cfg.Storage.Type = cnst.StorageTypeMemory
	// cfg.Storage.Type = cnst.StorageTypePostgres
	cfg.Storage.Conn = "postgres://calendar:CalendarPswd@127.0.0.1:5417/calendar_db?sslmode=disable"

	ctx := context.Background()
	log := logger.New(cfg.Logger.Level)
	store, err := storage.New(ctx, cfg, log)
	if err != nil {
		panic(err)
	}
	app := app.New(log, store)
	srv := httptest.NewServer(New(cfg, log, app).initRoute())

	now := time.Now()
	eventsExpected := testData(now)

	return &testHTTP{
		ctx:            context.Background(),
		store:          store,
		srv:            srv,
		client:         http.Client{Timeout: 30 * time.Second},
		now:            now,
		eventsExpected: eventsExpected,
		updateID:       3,
		expectedTitle:  "TestEventLastMonthUpdate",
	}
}

var s = &testHTTP{}

func init() {
	s = testInit()
}

func TestHTTPServerHelloWorld(t *testing.T) {
	url := s.srv.URL + "/hello"
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHTTPServerCreateEvent(t *testing.T) {
	for i, event := range s.eventsExpected {
		data, err := json.Marshal(event)
		require.NoError(t, err)

		url := s.srv.URL + "/event/"
		req, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewReader(data))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		b, err := io.ReadAll(res.Body)
		res.Body.Close()
		require.NoError(t, err)

		id := new(IDResponse)
		err = json.Unmarshal(b, id)
		require.NoError(t, err)
		require.Equal(t, i+1, id.Int())
	}
}

func TestHTTPServerGetEvents(t *testing.T) {
	url := s.srv.URL + "/event/"
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	err = json.Unmarshal(b, &s.eventsAll)
	require.NoError(t, err)
	require.Equal(t, 9, len(s.eventsAll))
}

func TestHTTPServerGetEventsByDay(t *testing.T) {
	url := s.srv.URL + "/event/day/"
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	events := new([]model.Event)
	err = json.Unmarshal(b, events)
	require.NoError(t, err)
	require.Equal(t, 1, len(*events))
}

func TestHTTPServerGetEventsByWeek(t *testing.T) {
	url := s.srv.URL + "/event/week/"
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	var events []model.Event
	err = json.Unmarshal(b, &events)
	require.NoError(t, err)
	require.Equal(t, 2, len(events))
}

func TestHTTPServerGetEventsByMonth(t *testing.T) {
	url := s.srv.URL + "/event/month/"
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	var events []model.Event
	err = json.Unmarshal(b, &events)
	require.NoError(t, err)
	require.Equal(t, 5, len(events))
}

func TestHTTPServerUpdateEvent(t *testing.T) {
	eventUpdate := s.eventsAll[s.updateID-1]
	eventUpdate.Title = s.expectedTitle

	data, err := json.Marshal(eventUpdate)
	require.NoError(t, err)

	url := s.srv.URL + "/event/"
	req, err := http.NewRequestWithContext(s.ctx, "PUT", url, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err = s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	var events []model.Event
	err = json.Unmarshal(b, &events)
	require.NoError(t, err)
	require.Equal(t, s.expectedTitle, events[s.updateID-1].Title)
}

func TestHTTPServerDeleteEvent(t *testing.T) {
	url := fmt.Sprintf("%s/event/%d/", s.srv.URL, s.updateID)

	req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	url = s.srv.URL + "/event/"
	req, err = http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err = s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	var events []model.Event
	err = json.Unmarshal(b, &events)
	require.NoError(t, err)
	require.Equal(t, len(s.eventsExpected)-1, len(events))
}

func TestHTTPServerDeleteOldEvents(t *testing.T) {
	url := s.srv.URL + "/event/old/"

	req, err := http.NewRequestWithContext(s.ctx, "DELETE", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := s.client.Do(req)
	res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	url = s.srv.URL + "/event/"
	req, err = http.NewRequestWithContext(s.ctx, "GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err = s.client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	var events []model.Event
	err = json.Unmarshal(b, &events)
	require.NoError(t, err)
	require.Equal(t, len(s.eventsExpected)-3, len(events))
}

func TestHTTPServerShutdown(t *testing.T) {
	s.srv.Close()
	err := s.store.Close()
	require.NoError(t, err)
}
