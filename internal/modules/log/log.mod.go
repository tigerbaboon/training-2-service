package log

import "app/internal/modules/config"

type LogModule struct {
	Svc *LogService
	Mid *LogMiddleware
}

func New(conf *config.ConfigService) *LogModule {
	svc := newService(conf.App())
	mid := NewMiddleware(conf, svc)
	return &LogModule{
		Svc: svc,
		Mid: mid,
	}
}
