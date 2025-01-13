package amenity

import "github.com/uptrace/bun"

type AmenityModule struct {
	Ctl *AmenityController
	Svc *AmenityService
}

func New(db *bun.DB) *AmenityModule {
	svc := newService(db)
	return &AmenityModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
