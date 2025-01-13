package zone

import (
	"app/app/modules/image"

	"github.com/uptrace/bun"
)

type ZoneModule struct {
	Ctl *ZoneController
	Svc *ZoneService
}

func New(db *bun.DB, ImageService *image.ImageService) *ZoneModule {
	svc := newService(db, ImageService)
	return &ZoneModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
