package activitylogs

import "github.com/uptrace/bun"

type ActivitylogsModule struct {
	Ctl *ActivitylogsController
	Svc *ActivitylogsService
}

func New(db *bun.DB) *ActivitylogsModule {
	svc := newService(db)
	return &ActivitylogsModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
