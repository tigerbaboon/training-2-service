package zone

import (
	"app/app/modules/base"
	zonedto "app/app/modules/zone/dto"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type ZoneController struct {
	zoneSvc *ZoneService
}

func newController(zoneService *ZoneService) *ZoneController {
	return &ZoneController{
		zoneSvc: zoneService,
	}
}

func (ctl *ZoneController) CreateZone(ctx *gin.Context) {
	var req zonedto.ZoneRequest

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	files := form.File["image"]
	if len(files) == 0 {
		base.BadRequest(ctx, "No files uploaded", nil)
		return
	}

	err = ctl.zoneSvc.CreateZone(ctx, &req, files)
	if err != nil {
		base.BadRequest(ctx, "Error creating zone: "+err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *ZoneController) GetZoneByID(ctx *gin.Context) {
	zone, err := ctl.zoneSvc.GetZoneByID(ctx, ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}
	base.Success(ctx, zone)
}

func (ctl *ZoneController) GetAllZone(ctx *gin.Context) {
	req := zonedto.ZoneGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	zoneResponse, count, err := ctl.zoneSvc.GetAllZone(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, zoneResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *ZoneController) UpdateZone(ctx *gin.Context) {
	zoneID := ctx.Param("id")
	req := zonedto.ZoneUpdateRequest{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	var files []*multipart.FileHeader
	if form != nil {
		files = form.File["images"]
	}

	err = ctl.zoneSvc.UpdateZone(ctx, zoneID, &req, files)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (z *ZoneController) DeleteZone(ctx *gin.Context) {
	err := z.zoneSvc.DeleteZone(ctx, ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}
	base.Success(ctx, "success")
}
