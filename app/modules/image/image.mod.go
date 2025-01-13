package image

import "github.com/uptrace/bun"

type ImageModule struct {
	Ctl *ImageController
	Svc *ImageService
}

func New(db *bun.DB) *ImageModule {
	svc := newService(db)
	return &ImageModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
