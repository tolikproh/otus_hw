package internalhttp

import (
	"net"
	"net/http"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
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

func loggingMiddleware(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		rw := newResWriter(w, log)
		next.ServeHTTP(rw, r)

		addr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			addr = "unknown"
		}
		log.Debug("http request",
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
