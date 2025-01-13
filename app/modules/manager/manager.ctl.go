package manager

import (
	"app/app/modules/base"
	managerdto "app/app/modules/manager/dto"

	"github.com/gin-gonic/gin"
)

type ManagerController struct {
	managerSvc *ManagerService
}

func newController(managerService *ManagerService) *ManagerController {
	return &ManagerController{
		managerSvc: managerService,
	}
}

func (*ManagerController) Get(ctx *gin.Context) {
	base.Success(ctx, "ok")
}

func (ctl *ManagerController) CreateManager(ctx *gin.Context) {
	var req managerdto.ManagerDTORequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	managerResponse, err := ctl.managerSvc.CreateManager(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, managerResponse)
}

func (ctl *ManagerController) UpdateManager(ctx *gin.Context) {
	username := ctx.Param("username")

	var req managerdto.ManagerUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	managerResponse, err := ctl.managerSvc.UpdateManager(ctx, username, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, managerResponse)
}

func (ctl *ManagerController) LoginManager(ctx *gin.Context) {
	var reqlogin managerdto.ManagerLoginRequest
	if err := ctx.ShouldBindJSON(&reqlogin); err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	managerResponse, err := ctl.managerSvc.LoginManager(ctx, &reqlogin)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, managerResponse)
}
