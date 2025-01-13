package routes

import (
	"app/app/modules"
	"app/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WarpH(router *gin.RouterGroup, prefix string, handler http.Handler) {
	router.Any(fmt.Sprintf("%s/*w", prefix), gin.WrapH(http.StripPrefix(fmt.Sprintf("%s%s", router.BasePath(), prefix), handler)))
}

func api(r *gin.RouterGroup, mod *modules.Modules) {

	md := middleware.CheckJwtAuth()
	log := middleware.NewLogResponse()

	// User
	r.GET("/user/show", md, mod.User.Ctl.GetUser)
	r.POST("/user/create", mod.User.Ctl.CreateUser)
	r.POST("/user/login", mod.User.Ctl.Login)
	r.PATCH("/user/update/:username", mod.User.Ctl.UpdateUser)

	// Promote
	r.GET("/promote/show", mod.Promote.Ctl.GetAllPromotes)
	r.GET("/promote/show/:id", mod.Promote.Ctl.GetPromoteByID)

	// Travel
	r.GET("/travel/show", mod.Travel.Ctl.GetAllTravels)
	r.GET("/travel/show/:id", mod.Travel.Ctl.GetTravelByID)

	// Zone
	r.GET("/zone/show", mod.Zone.Ctl.GetAllZone)
	r.GET("/zone/show/:id", mod.Zone.Ctl.GetZoneByID)

	// Amenity
	r.GET("/amenity/show", mod.Amenity.Ctl.GetAllAmenity)
	r.GET("/amenity/show/:id", mod.Amenity.Ctl.GetAmenityByID)

	// House
	r.GET("/house/show/:id", mod.House.Ctl.GetHouseByID)
	r.GET("/house/show", mod.House.Ctl.GetAllHouses)
	r.GET("/house/get", md, mod.House.Ctl.GetHousesByProfile)
	r.GET("/house/getnearby", mod.House.Ctl.GetNearbyHouses)
	r.GET("/house/getprice", mod.House.Ctl.GetPriceRange)
	r.POST("/house/create", md, log, mod.House.Ctl.CreateHouse)
	r.PATCH("/house/updatestatus/:id", md, log, mod.House.Ctl.UpdateStatusHouse)
	r.PATCH("/house/update/:id", md, log, mod.House.Ctl.UpdateHouse)
	r.DELETE("/house/delete/:id", md, log, mod.House.Ctl.DeleteHouse)

}

func apiadmin(r *gin.RouterGroup, mod *modules.Modules) {

	md := middleware.CheckJwtAuth()
	log := middleware.NewLogResponse()

	// Menager
	r.POST("/manager/create", log, mod.Manager.Ctl.CreateManager)
	r.POST("/manager/login", mod.Manager.Ctl.LoginManager)
	r.PATCH("/manager/update/:username", md, log, mod.Manager.Ctl.UpdateManager)

	// Promote
	r.GET("/promote/show", md, mod.Promote.Ctl.GetAllPromotes)
	r.GET("/promote/show/:id", md, mod.Promote.Ctl.GetPromoteByID)
	r.POST("/promote/create", md, log, mod.Promote.Ctl.CreatePromote)
	r.PATCH("/promote/update/:id", md, log, mod.Promote.Ctl.UpdatePromote)
	r.PATCH("/promote/update/status/:id", md, log, mod.Promote.Ctl.UpdatePromoteStatus)
	r.DELETE("/promote/delete/:id", md, log, mod.Promote.Ctl.DeletePromote)

	// Travel
	r.GET("/travel/show", md, mod.Travel.Ctl.GetAllTravelForAdmin)
	r.GET("/travel/show/:id", md, mod.Travel.Ctl.GetTravelByID)
	r.POST("/travel/create", md, log, mod.Travel.Ctl.CreateTravel)
	r.PATCH("/travel/update/:id", md, log, mod.Travel.Ctl.UpdateTravel)
	r.PATCH("/travel/update/status/:id", md, log, mod.Travel.Ctl.UpdateTravelStatus)
	r.DELETE("/travel/delete/:id", md, log, mod.Travel.Ctl.DeleteTravel)

	// Zone
	r.GET("/zone/show", md, mod.Zone.Ctl.GetAllZone)
	r.GET("/zone/show/:id", md, mod.Zone.Ctl.GetZoneByID)
	r.POST("/zone/create", md, log, mod.Zone.Ctl.CreateZone)
	r.PATCH("zone/update/:id", md, log, mod.Zone.Ctl.UpdateZone)
	r.DELETE("/zone/delete/:id", md, log, mod.Zone.Ctl.DeleteZone)

	// House
	r.GET("/house/show/:id", md, mod.House.Ctl.GetHouseByID)
	r.GET("/house/show", md, mod.House.Ctl.GetAllHousesForAdmin)
	r.GET("/house/gethistory", md, mod.House.Ctl.GetHouseHistory)
	r.GET("/house/getcount", md, mod.House.Ctl.GetHouseCountByZone)
	r.GET("/house/getconfirm", md, mod.House.Ctl.GetHousesConfirmation)
	r.POST("/house/create", md, log, mod.House.Ctl.CreateHouse)
	r.PATCH("/house/update/:id", md, log, mod.House.Ctl.UpdateHouse)
	r.PATCH("/house/update/rec/:id", md, log, mod.House.Ctl.UpdateRecommendHouse)
	r.PATCH("/house/update/confirm/:id", md, log, mod.House.Ctl.UpdateConfirmation)
	r.DELETE("/house/delete/:id", md, log, mod.House.Ctl.DeleteHouse)

	// Amenity
	r.POST("/amenity/create", md, log, mod.Amenity.Ctl.CreateAmenity)
	r.GET("/amenity/show", md, mod.Amenity.Ctl.GetAllAmenity)
	r.GET("/amenity/show/:id", md, mod.Amenity.Ctl.GetAmenityByID)
	r.PATCH("/amenity/update/:id", md, log, mod.Amenity.Ctl.UpdateAmenity)
	r.DELETE("/amenity/delete/:id", md, log, mod.Amenity.Ctl.DeleteAmenity)
}
