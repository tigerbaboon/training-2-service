package collector

import (
	"app/config"
	"app/internal/service/provider"
	"context"
)

type OTELCollectorModule struct {
	Service *OTELCollectorService
}

var _ provider.Close = (*OTELCollectorModule)(nil)

func New(conf *config.Config) *OTELCollectorModule {
	return &OTELCollectorModule{
		Service: newService(conf),
	}
}

func (m *OTELCollectorModule) Close(ctx context.Context) error {
	if m.Service == nil {
		return nil
	}
	return m.Service.close(ctx)
}
