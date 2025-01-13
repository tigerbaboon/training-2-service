package database

import (
	"app/internal/modules/config"
	"app/internal/service/provider"
	"context"
)

type DatabaseModule struct {
	Svc *DatabaseService
}

var _ provider.Close = (*DatabaseModule)(nil)

func New(conf *config.ConfigService) *DatabaseModule {
	dbOpts := conf.Database().Sql
	service := newService(dbOpts)
	return &DatabaseModule{
		Svc: service,
	}
}

func (db *DatabaseModule) Close(ctx context.Context) error {
	return db.Svc.close(ctx)
}
