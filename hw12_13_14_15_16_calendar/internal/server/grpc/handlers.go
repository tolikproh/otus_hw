package grpc

import (
	"context"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/pb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type eventDataHandler struct {
	pb.UnimplementedCalendarServer
	log *logger.Logger
	app *app.App
}

func newEventDataHandler(log *logger.Logger, app *app.App) *eventDataHandler {
	return &eventDataHandler{log: log, app: app}
}

//nolint:gofumpt
func (h *eventDataHandler) GetEvents(ctx context.Context,
	params *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := h.app.GetEvents(ctx)
	if err != nil {
		return nil, err
	}
	return modelEvensToResponse(events), nil
}

//nolint:gofumpt
func (h *eventDataHandler) GetEventsByDay(ctx context.Context,
	params *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := h.app.GetEventsByDay(ctx)
	if err != nil {
		return nil, err
	}

	return modelEvensToResponse(events), nil
}

//nolint:gofumpt
func (h *eventDataHandler) GetEventsByWeek(ctx context.Context,
	params *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := h.app.GetEventsByWeek(ctx)
	if err != nil {
		return nil, err
	}

	return modelEvensToResponse(events), nil
}

//nolint:gofumpt
func (h *eventDataHandler) GetEventsByMonth(ctx context.Context,
	params *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := h.app.GetEventsByMonth(ctx)
	if err != nil {
		return nil, err
	}

	return modelEvensToResponse(events), nil
}

//nolint:gofumpt
func (h *eventDataHandler) CreateEvent(ctx context.Context,
	event *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	newEvent := pbEventToModelEvent(event.Event.Event)
	id, err := h.app.CreateEvent(ctx, newEvent)
	if err != nil {
		return nil, err
	}
	eventID := new(pb.CreateEventResponse)
	eventID.Id = int32(id) //nolint:gosec

	return eventID, nil
}

//nolint:gofumpt
func (h *eventDataHandler) UpdateEvent(ctx context.Context,
	event *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	newEvent := pbEventToModelEvent(event.Event.Event)
	err := h.app.UpdateEvent(ctx, newEvent.ID, newEvent)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//nolint:gofumpt
func (h *eventDataHandler) DeleteEvent(ctx context.Context,
	params *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := h.app.DeleteEvent(ctx, int(params.GetId()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//nolint:gofumpt
func (h *eventDataHandler) DeleteEventsOld(ctx context.Context,
	params *pb.DeleteEventsOldRequest) (*pb.DeleteEventsOldResponse, error) {
	err := h.app.DeleteEventsOld(ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func modelEventToCreateEventRequest(event model.Event) *pb.CreateEventRequest {
	res := new(pb.CreateEventRequest)
	res.Event = pbEventToEventRequest(modelEventToPbEvent(event))

	return res
}

func modelEventToUpdateEventRequest(event model.Event) *pb.UpdateEventRequest {
	res := new(pb.UpdateEventRequest)
	res.Event = pbEventToEventRequest(modelEventToPbEvent(event))

	return res
}

func pbEventToEventRequest(event *pb.Event) *pb.EventRequest {
	res := new(pb.EventRequest)
	res.Event = event
	return res
}

func modelEvensToResponse(events []model.Event) *pb.GetEventsResponse {
	res := new(pb.GetEventsResponse)
	for _, event := range events {
		res.Events = append(res.Events, modelEventToPbEvent(event))
	}
	return res
}

func modelEventToPbEvent(event model.Event) *pb.Event {
	tsDateTimeStart := timestamppb.New(event.DateTimeStart)
	tsDateTimeEnd := timestamppb.New(event.DateTimeEnd)
	notificationDuration := durationpb.New(*event.NotificationDuration)
	res := new(pb.Event)
	res.Id = int32(event.ID) //nolint:gosec
	res.Title = event.Title
	res.DateTimeStart = tsDateTimeStart
	res.DateTimeEnd = tsDateTimeEnd
	res.Description = event.Description
	res.UserId = int32(event.UserID) //nolint:gosec
	res.NotificationDuration = notificationDuration

	return res
}

func pbEventToModelEvent(event *pb.Event) model.Event {
	notificationDuration := event.GetNotificationDuration().AsDuration()
	res := model.Event{}
	res.ID = int(event.GetId())
	res.Title = event.GetTitle()
	res.DateTimeStart = event.GetDateTimeStart().AsTime()
	res.DateTimeEnd = event.GetDateTimeEnd().AsTime()
	res.Description = event.GetDescription()
	res.UserID = int(event.GetUserId())
	res.NotificationDuration = &notificationDuration

	return res
}
