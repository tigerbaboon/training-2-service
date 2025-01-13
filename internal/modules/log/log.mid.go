package log

import "app/internal/modules/config"

type LogMiddleware struct {
	Config *config.ConfigService
	Svc    *LogService
}

func NewMiddleware(conf *config.ConfigService, svc *LogService) *LogMiddleware {
	return &LogMiddleware{
		Config: conf,
		Svc:    svc,
	}
}
