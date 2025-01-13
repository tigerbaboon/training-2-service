package amenity

import (
	amenitydto "app/app/modules/amenity/dto"
	"app/app/modules/base"

	"github.com/gin-gonic/gin"
)

type AmenityController struct {
	amenitySvc *AmenityService
}

func newController(amenityService *AmenityService) *AmenityController {
	return &AmenityController{
		amenitySvc: amenityService,
	}
}

func (ctl *AmenityController) CreateAmenity(ctx *gin.Context) {
	var req amenitydto.AmenityRequest

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	err := ctl.amenitySvc.CreateAmenity(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}
	base.Success(ctx, "success")
}

func (ctl *AmenityController) GetAmenityByID(ctx *gin.Context) {
	amenity, err := ctl.amenitySvc.GetAmenityByID(ctx, ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}
	base.Success(ctx, amenity)
}

func (ctl *AmenityController) GetAllAmenity(ctx *gin.Context) {
	req := amenitydto.AmenityGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	amenityResponse, count, err := ctl.amenitySvc.GetAllAmenity(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, amenityResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *AmenityController) UpdateAmenity(ctx *gin.Context) {
	amntID := ctx.Param("id")
	req := amenitydto.AmenityRequest{}
	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	amenityResponse, err := ctl.amenitySvc.UpdateAmenity(ctx, amntID, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, amenityResponse)
}

func (ctl *AmenityController) DeleteAmenity(ctx *gin.Context) {
	err := ctl.amenitySvc.DeleteAmenity(ctx, ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}
	base.Success(ctx, nil)
}
