package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/model"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/pkg/httperr"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func (s *Server) getEventsHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	events, err := s.app.GetEvents(ctx)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return events, nil
}

func (s *Server) getEventsByDayHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	events, err := s.app.GetEventsByDay(ctx)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return events, nil
}

func (s *Server) getEventsByWeekHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	events, err := s.app.GetEventsByWeek(ctx)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return events, nil
}

func (s *Server) getEventsByMonthHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	events, err := s.app.GetEventsByMonth(ctx)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}
	return events, nil
}

func (s *Server) createEventHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var event model.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	id, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	event.ID = id

	return ToIDResponse(id), nil
}

func (s *Server) updateEventHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	var event model.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err = s.app.UpdateEvent(ctx, event.ID, event)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) deleteEventHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusBadRequest)
	}

	err = s.app.DeleteEvent(ctx, id)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *Server) deleteEventsOldHandler(ctx context.Context, r *http.Request) (interface{}, error) {
	err := s.app.DeleteEventsOld(ctx)
	if err != nil {
		return nil, httperr.WithHTTPStatus(err, http.StatusInternalServerError)
	}

	return nil, nil
}
