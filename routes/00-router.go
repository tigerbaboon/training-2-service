package routes

import (
	"app/app/modules"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router Router to fiber framework
func Router(app *gin.Engine, mod *modules.Modules) {

	app.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
		return
	})
	app.Use(otelgin.Middleware(mod.Conf.Svc.App().AppName))
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"*"},
		AllowHeaders:           []string{"*"},
		AllowCredentials:       true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             false,
	}))

	web(app.Group("/"), mod)
	api(app.Group("/api/user"), mod)
	apiadmin(app.Group("/api/admin"), mod)
}
