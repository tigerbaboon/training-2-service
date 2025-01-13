package travel

import (
	"app/app/modules/base"
	traveldto "app/app/modules/travel/dto"
	"app/internal/modules/log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TravelController struct {
	travelSvc *TravelService
}

func newController(travelService *TravelService) *TravelController {
	return &TravelController{
		travelSvc: travelService,
	}
}

func (*TravelController) Get(ctx *gin.Context) {
	base.Success(ctx, "ok")
}

func (ctl *TravelController) CreateTravel(ctx *gin.Context) {
	req := traveldto.TravelRequest{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	imageMainFile, err := ctx.FormFile("image_main")
	if err != nil || imageMainFile == nil {
		base.BadRequest(ctx, "Main image is required", nil)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	imageFiles := form.File["image"]
	if len(imageFiles) == 0 {
		base.BadRequest(ctx, "No files uploaded", nil)
		return
	}

	err = ctl.travelSvc.CreateTravel(ctx, &req, imageMainFile, imageFiles)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *TravelController) GetTravelByID(ctx *gin.Context) {
	travelID := ctx.Param("id")

	travelResponse, err := ctl.travelSvc.GetTravelByID(ctx, travelID)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, travelResponse)
}

func (ctl *TravelController) GetAllTravels(ctx *gin.Context) {
	req := traveldto.TravelGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	travelResponse, count, err := ctl.travelSvc.GetAllTravels(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, travelResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}
func (ctl *TravelController) GetAllTravelForAdmin(ctx *gin.Context) {
	req := traveldto.TravelGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		log.Info("error : %v", err)
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	travelResponse, count, err := ctl.travelSvc.GetAllTravelForAdmin(ctx, &req)
	if err != nil {
		log.Info("error : %v", err)
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, travelResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *TravelController) UpdateTravel(ctx *gin.Context) {
	travelID := ctx.Param("id")
	req := traveldto.TravelUpdateRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	var ImageMainFile *multipart.FileHeader
	var imageFiles []*multipart.FileHeader

	if form != nil {
		ImageMainFiles := form.File["image_main"]
		if len(ImageMainFiles) > 0 {
			ImageMainFile = ImageMainFiles[0]
		}

		imageFiles = form.File["image"]
	}

	err = ctl.travelSvc.UpdateTravel(ctx, travelID, &req, ImageMainFile, imageFiles)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *TravelController) UpdateTravelStatus(ctx *gin.Context) {
	travelID := ctx.Param("id")
	req := traveldto.TravelUpdateStatus{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	travelResponse, err := ctl.travelSvc.UpdateTravelStatus(ctx, travelID, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, travelResponse)
}

func (ctl *TravelController) DeleteTravel(ctx *gin.Context) {
	travelID := ctx.Param("id")

	err := ctl.travelSvc.DeleteTravel(ctx, travelID)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "Travel deleted")
}
