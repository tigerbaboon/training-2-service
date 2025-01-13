package travel

import (
	"app/app/modules/image"

	"github.com/uptrace/bun"
)

type TravelModule struct {
	Ctl *TravelController
	Svc *TravelService
}

func New(db *bun.DB, ImageService *image.ImageService) *TravelModule {
	svc := newService(db, ImageService)
	return &TravelModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
