package house

import (
	"app/app/modules/amenity"
	"app/app/modules/contact"
	"app/app/modules/image"
	"app/app/modules/manager"
	"app/app/modules/user"
	"app/app/modules/zone"

	"github.com/uptrace/bun"
)

type HouseModule struct {
	Ctl *HouseController
	Svc *HouseService
}

func New(db *bun.DB, ImageService *image.ImageService, ContactService *contact.ContactService, AmenityService *amenity.AmenityService, ZoneService *zone.ZoneService, UserService *user.UserService, ManagerService *manager.ManagerService) *HouseModule {
	svc := newService(db, ImageService, ContactService, AmenityService, ZoneService, UserService, ManagerService)
	return &HouseModule{
		Ctl: newController(svc),
		Svc: svc,
	}
}
