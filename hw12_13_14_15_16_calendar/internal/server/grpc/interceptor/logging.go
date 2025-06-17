package interceptor

import (
	"context"
	"strings"
	"time"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Interceptor struct {
	log *logger.Logger
}

func New(log *logger.Logger) *Interceptor {
	return &Interceptor{log}
}

//nolint:gofumpt
func (i *Interceptor) Logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		result, err := handler(ctx, req)

		i.log.Info("grpc request",
			"start time", startTime.UTC(),
			"method", info.FullMethod[strings.LastIndex(info.FullMethod, "/"):],
			"latency [ms]", time.Since(startTime).Microseconds(),
			"user agent", getUserAgent(ctx))

		return result, err
	}
}

func getUserAgent(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	agents, ok := md["user-agent"]
	if !ok || len(agents) == 0 {
		return ""
	}

	return agents[0]
}
