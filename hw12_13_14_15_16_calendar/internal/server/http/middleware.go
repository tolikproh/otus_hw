package internalhttp

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/pkg/httperr"
)

type resWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
	log        *logger.Logger
}

func newResWriter(w http.ResponseWriter, log *logger.Logger) *resWriter {
	return &resWriter{w, http.StatusOK, 0, log}
}

func (rw *resWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *resWriter) Write(bytes []byte) (int, error) {
	l, err := rw.ResponseWriter.Write(bytes)
	if err != nil {
		rw.log.Error("write response error", "error", err)
		return 0, err
	}
	return l, nil
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		rw := newResWriter(w, s.log)
		next.ServeHTTP(rw, r)

		addr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			addr = "unknown"
		}
		s.log.Info("http request",
			"address", addr,
			"start time", startTime.UTC(),
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"status code", rw.statusCode,
			"latency [ms]", time.Since(startTime).Microseconds(),
			"user agent", r.UserAgent())
	})
}

type errorHTTP struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Error  string `json:"error"`
}

type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

func (s *Server) serveHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		response := new(http.Response)

		ctx, cansel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cansel()

		data, err := h(ctx, r)
		if err != nil {
			code := httperr.HTTPStatus(err)

			response.StatusCode = code
			r.Response = response

			errorHTTP := new(errorHTTP)
			errorHTTP.Method = r.Method
			errorHTTP.Path = r.RequestURI
			errorHTTP.Error = err.Error()

			s.log.Error("http error request",
				"method", errorHTTP.Method,
				"path", errorHTTP.Path,
				"status code", code,
				"error", errorHTTP.Error,
				"user agent", r.UserAgent())

			b, _ := json.Marshal(errorHTTP)
			w.WriteHeader(code)
			w.Write(b)

			return
		}

		response.StatusCode = http.StatusOK
		r.Response = response
		b, _ := json.Marshal(data)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
