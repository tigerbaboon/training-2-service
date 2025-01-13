package manager

import "github.com/uptrace/bun"

type ManagerModule struct {
	Ctl *ManagerController
	Svc *ManagerService
}

func New(db *bun.DB) *ManagerModule {
	svc := newService(db)
	return &ManagerModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
