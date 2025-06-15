package internalhttp

import "net/http"

func (s *Server) initRoute() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/hello/", http.HandlerFunc(helloHandler))

	mux.Handle("GET /event/", s.serveHandler(s.getEventsHandler))
	mux.Handle("GET /event/day/", s.serveHandler(s.getEventsByDayHandler))
	mux.Handle("GET /event/week/", s.serveHandler(s.getEventsByWeekHandler))
	mux.Handle("GET /event/month/", s.serveHandler(s.getEventsByMonthHandler))
	mux.Handle("POST /event/", s.serveHandler(s.createEventHandler))
	mux.Handle("PUT /event/", s.serveHandler(s.updateEventHandler))
	mux.Handle("DELETE /event/{id}/", s.serveHandler(s.deleteEventHandler))
	mux.Handle("DELETE /event/old/", s.serveHandler(s.deleteEventsOldHandler))

	var handler http.Handler = mux
	handler = s.loggingMiddleware(handler)

	return handler
}
