package redis

import (
	"app/internal/modules/config"
	"app/internal/service/provider"
	"context"
)

type RedisModule struct {
	Svc *RedisService
}

var _ provider.Close = (*RedisModule)(nil)

func New(conf *config.ConfigService) *RedisModule {
	svc := newService(conf)
	return &RedisModule{
		Svc: svc,
	}
}

func (db *RedisModule) Close(ctx context.Context) error {
	return db.Svc.close(ctx)
}
