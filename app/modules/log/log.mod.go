package log

import "github.com/uptrace/bun"

type LogModule struct {
	Ctl *LogController
	Svc *LogService
}

func New(db *bun.DB) *LogModule {
	svc := newService(db)
	return &LogModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
