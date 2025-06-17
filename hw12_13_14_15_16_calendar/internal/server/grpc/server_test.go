package grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/config"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/pb"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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

type testGRPC struct {
	ctx            context.Context
	cansel         context.CancelFunc
	store          storage.Storage
	conn           *grpc.ClientConn
	client         pb.CalendarClient
	server         *grpc.Server
	now            time.Time
	eventsExpected []model.Event
	updateID       int
	eventsAll      []model.Event
	expectedTitle  string
}

func testInit() *testGRPC {
	ctx, cansel := context.WithCancel(context.Background())

	cfg := new(config.Config)
	cfg.Logger.Level = cnst.LoggerLevelDebug
	cfg.Storage.Type = cnst.StorageTypeMemory
	// cfg.Storage.Type = cnst.StorageTypePostgres
	cfg.Storage.Conn = "postgres://calendar:CalendarPswd@127.0.0.1:5417/calendar_db?sslmode=disable"

	log := logger.New(cfg.Logger.Level)
	store, err := storage.New(ctx, cfg, log)
	if err != nil {
		panic(err)
	}
	app := app.New(log, store)

	server := grpc.NewServer()
	pb.RegisterCalendarServer(server, newEventDataHandler(log, app))

	listener := bufconn.Listen(1024 * 1024)
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			err := server.Serve(listener)
			if err != nil {
				panic(err)
			}
		}
	}()

	connFunc := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(connFunc),
	}
	//nolint:staticcheck
	conn, err := grpc.DialContext(ctx, "", opts...)
	if err != nil {
		panic(err)
	}

	client := pb.NewCalendarClient(conn)

	now := time.Now()
	eventsExpected := testData(now)

	return &testGRPC{
		ctx:            context.Background(),
		cansel:         cansel,
		store:          store,
		conn:           conn,
		client:         client,
		server:         server,
		now:            now,
		eventsExpected: eventsExpected,
		updateID:       3,
		expectedTitle:  "TestEventLastMonthUpdate",
	}
}

var s = &testGRPC{}

func init() {
	s = testInit()
}

func TestGRPCServerCreateEvent(t *testing.T) {
	for i, event := range s.eventsExpected {
		data := modelEventToCreateEventRequest(event)
		res, err := s.client.CreateEvent(s.ctx, data)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, i+1, int(res.GetId()))
	}
}

func TestGRPCServerGetEvents(t *testing.T) {
	res, err := s.client.GetEvents(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	for _, v := range res.Events {
		s.eventsAll = append(s.eventsAll, pbEventToModelEvent(v))
	}
	require.Equal(t, 9, len(s.eventsAll))
}

func TestGRPCServerGetEventsByDay(t *testing.T) {
	res, err := s.client.GetEventsByDay(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Events))
}

func TestGRPCServerGetEventsByWeek(t *testing.T) {
	res, err := s.client.GetEventsByWeek(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	require.Equal(t, 2, len(res.Events))
}

func TestGRPCServerGetEventsByMonth(t *testing.T) {
	res, err := s.client.GetEventsByMonth(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	require.Equal(t, 5, len(res.Events))
}

func TestGRPCServerUpdateEvent(t *testing.T) {
	eventUpdate := s.eventsAll[s.updateID-1]
	eventUpdate.Title = s.expectedTitle
	req := modelEventToUpdateEventRequest(eventUpdate)
	_, err := s.client.UpdateEvent(s.ctx, req)
	require.NoError(t, err)

	res, err := s.client.GetEvents(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	require.Equal(t, s.expectedTitle, res.Events[s.updateID-1].Title)
}

func TestGRPCServerDeleteEvent(t *testing.T) {
	req := new(pb.DeleteEventRequest)
	req.Id = int32(s.updateID)
	_, err := s.client.DeleteEvent(s.ctx, req)
	require.NoError(t, err)

	res, err := s.client.GetEvents(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	require.Equal(t, len(s.eventsExpected)-1, len(res.Events))
}

func TestGRPCServerDeleteOldEvents(t *testing.T) {
	_, err := s.client.DeleteEventsOld(s.ctx, &pb.DeleteEventsOldRequest{})
	require.NoError(t, err)

	res, err := s.client.GetEvents(s.ctx, &pb.GetEventsRequest{})
	require.NoError(t, err)
	require.Equal(t, len(s.eventsExpected)-3, len(res.Events))
}

func TestGRPCServerShutdown(t *testing.T) {
	s.cansel()
	err := s.conn.Close()
	require.NoError(t, err)
	s.server.GracefulStop()
}
