package user

import (
	"app/app/modules/base"
	userdto "app/app/modules/user/dto"

	// "app/app/provider/database"
	// "app/app/provider/redis"

	"github.com/gin-gonic/gin"
)

type UserCTL struct {
	userService *UserService
}

func newCTL(userService *UserService) *UserCTL {
	return &UserCTL{
		userService: userService,
	}
}

func (u *UserCTL) Get(ctx *gin.Context) {
	base.Success(ctx, userdto.UserResponse{})
}

func (ctl *UserCTL) CreateUser(ctx *gin.Context) {
	var user userdto.UserDTORequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		base.BadRequest(ctx, "Invalid request format!", nil)
		return
	}

	err := ctl.userService.CreateUser(ctx, &user)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *UserCTL) Login(ctx *gin.Context) {
	var reqlogin userdto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&reqlogin); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	userResponse, err := ctl.userService.Login(ctx, reqlogin)
	if err != nil {
		base.BadRequest(ctx, "Data not found", nil)
		return
	}

	base.Success(ctx, userResponse)
}

func (ctl *UserCTL) UpdateUser(ctx *gin.Context) {
	username := ctx.Param("username")

	var req userdto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	err := ctl.userService.UpdateUser(ctx, username, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *UserCTL) GetUser(ctx *gin.Context) {
	userx, exists := ctx.Get("userId")
	if !exists {
		base.BadRequest(ctx, "User not found", nil)
		return
	}

	userResponse, err := ctl.userService.GetUser(ctx, userx.(string))
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, userResponse)
}
