package log

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
)

type statusResponse struct {
	Status struct {
		Code int32
	}
}

func (l *LogMiddleware) UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	rt := time.Since(startTime)
	respStatus := &statusResponse{}
	copier.CopyWithOption(respStatus, resp, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})

	var log *slog.Logger
	if respStatus.Status.Code != 200 || l.Config.App().Debug || rt > 250*time.Millisecond {
		log = slog.With(
			slog.String("full_method", info.FullMethod),
			slog.Any("request", req),
			slog.Duration("response_time", rt),
		)
		if err == nil && respStatus.Status.Code == 200 {
			if l.Config.App().Debug {
				log = log.With(slog.Any("response", resp))
			}
			log.Info(fmt.Sprintf("Access | %20s | %s", rt, info.FullMethod))
		} else {
			log = log.With(slog.String("error", err.Error()), slog.Any("response", resp))
			log.Error(fmt.Sprintf("Error  | %20s | %s", rt, info.FullMethod))
		}
	}
	// latency, err := global.MeterProvider().Meter("grpc.go.log.v1").SyncFloat64().Counter("grpc.server.request.latency")
	return resp, err

}
