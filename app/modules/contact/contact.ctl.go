package contact

import (
	"app/app/modules/base"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	contactSvc *ContactService
}

func newController(contactService *ContactService) *ContactController {
	return &ContactController{
		contactSvc: contactService,
	}
}

func (*ContactController) Get(ctx *gin.Context) {
	base.Success(ctx, "ok")
}
