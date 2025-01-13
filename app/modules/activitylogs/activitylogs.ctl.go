package activitylogs

import (
	"app/app/modules/base"

	"github.com/gin-gonic/gin"
)

type ActivitylogsController struct {
	activitylogsSvc *ActivitylogsService
}

func newController(activitylogsSvc *ActivitylogsService) *ActivitylogsController {
	return &ActivitylogsController{
		activitylogsSvc: activitylogsSvc,
	}
}

func (*ActivitylogsController) Get(ctx *gin.Context) {
	base.Success(ctx, "ok")
}
