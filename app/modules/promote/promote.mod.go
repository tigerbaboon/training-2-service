package promote

import (
	"app/app/modules/image"

	"github.com/uptrace/bun"
)

type PromoteModule struct {
	Ctl *PromoteController
	Svc *PromoteService
}

func New(db *bun.DB, ImageService *image.ImageService) *PromoteModule {
	svc := newService(db, ImageService)
	return &PromoteModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
