package contact

import "github.com/uptrace/bun"

type ContactModule struct {
	Ctl *ContactController
	Svc *ContactService
}

func New(db *bun.DB) *ContactModule {
	svc := newService(db)
	return &ContactModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
