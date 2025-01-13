package modules

import (
	"app/app/modules/activitylogs"
	"app/app/modules/amenity"
	"app/app/modules/contact"
	"app/app/modules/example"
	"app/app/modules/house"
	"app/app/modules/image"
	"app/app/modules/manager"
	"app/app/modules/promote"
	"app/app/modules/travel"
	"app/app/modules/user"
	"app/app/modules/zone"
	"app/internal/modules/config"
	"app/internal/modules/database"
	"app/internal/modules/log"
	"app/internal/modules/otel/collector"
	"sync"
)

type Modules struct {
	Conf         *config.ConfigModule
	Log          *log.LogModule
	OTEL         *collector.OTELCollectorModule
	DB           *database.DatabaseModule
	Example      *example.ExampleModule
	User         *user.UserModule
	Manager      *manager.ManagerModule
	Amenity      *amenity.AmenityModule
	Promote      *promote.PromoteModule
	House        *house.HouseModule
	Image        *image.ImageModule
	Travel       *travel.TravelModule
	Zone         *zone.ZoneModule
	Contact      *contact.ContactModule
	Acticitylogs *activitylogs.ActivitylogsModule
}

func modulesInit() {
	conf := config.New()

	logMod := log.New(conf.Svc)

	otel := collector.New(conf.Svc.App())
	log.Info("otel module initialized")

	db := database.New(conf.Svc)
	log.Info("database module initialized")

	user := user.New(db.Svc.DB())
	log.Info("user module initialized")

	manager := manager.New(db.Svc.DB())
	log.Info("manager module initialized")

	contact := contact.New(db.Svc.DB())
	log.Info("contact module initialized")

	image := image.New(db.Svc.DB())
	log.Info("image module initialized")

	zone := zone.New(db.Svc.DB(), image.Svc)
	log.Info("zone module initialized")

	amenity := amenity.New(db.Svc.DB())
	log.Info("amenity module initialized")

	house := house.New(db.Svc.DB(), image.Svc, contact.Svc, amenity.Svc, zone.Svc, user.Service, manager.Svc)
	log.Info("house module initialized")

	travel := travel.New(db.Svc.DB(), image.Svc)
	log.Info("travel module initialized")

	example := example.New(user.Service)
	log.Info("example module initialized")

	promote := promote.New(db.Svc.DB(), image.Svc)
	log.Info("promote module initialized")

	activitylogs := activitylogs.New(db.Svc.DB())
	log.Info("activitylogs module initialized")

	mod = &Modules{
		Conf:         conf,
		Log:          logMod,
		OTEL:         otel,
		DB:           db,
		User:         user,
		Manager:      manager,
		House:        house,
		Image:        image,
		Travel:       travel,
		Contact:      contact,
		Example:      example,
		Amenity:      amenity,
		Promote:      promote,
		Zone:         zone,
		Acticitylogs: activitylogs,
	}
	log.Info("all modules initialized")
}

var (
	once sync.Once
	mod  *Modules
)

func Get() *Modules {
	once.Do(modulesInit)

	return mod
}
