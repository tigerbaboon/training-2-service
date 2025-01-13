package promote

import (
	"app/app/modules/base"
	promotedto "app/app/modules/promote/dto"
	"app/internal/modules/log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type PromoteController struct {
	promoteSvc *PromoteService
}

func newController(promoteService *PromoteService) *PromoteController {
	return &PromoteController{
		promoteSvc: promoteService,
	}
}

func (ctl *PromoteController) CreatePromote(ctx *gin.Context) {
	req := promotedto.PromoteDTORequest{}

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

	promoteResponse, err := ctl.promoteSvc.CreatePromote(ctx, &req, files)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, promoteResponse)
}

func (ctl *PromoteController) GetPromoteByID(ctx *gin.Context) {
	promoteID := ctx.Param("id")

	promoteResponse, err := ctl.promoteSvc.GetPromoteByID(ctx, promoteID)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, promoteResponse)
}

func (ctl *PromoteController) GetAllPromotes(ctx *gin.Context) {
	req := promotedto.PromoteGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	promoteResponse, count, err := ctl.promoteSvc.GetAllPromotes(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, promoteResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *PromoteController) UpdatePromote(ctx *gin.Context) {
	promoteID := ctx.Param("id")

	req := promotedto.PromoteUpdateRequest{}
	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		log.Info("Error : %v", err)
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	var imageFile *multipart.FileHeader
	if form != nil {
		files := form.File["image"]
		if len(files) > 0 {
			imageFile = files[0]
		}
	}

	log.Info("Image file: %v", imageFile)

	err = ctl.promoteSvc.UpdatePromote(ctx, promoteID, &req, imageFile)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *PromoteController) UpdatePromoteStatus(ctx *gin.Context) {
	promoteID := ctx.Param("id")
	req := promotedto.PromoteUpdateStatus{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	promoteResponse, err := ctl.promoteSvc.UpdatePromoteStatus(ctx, promoteID, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, promoteResponse)
}

func (ctl *PromoteController) DeletePromote(ctx *gin.Context) {
	promoteID := ctx.Param("id")

	err := ctl.promoteSvc.DeletePromote(ctx, promoteID)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "Promote deleted")
}
