package house

import (
	"app/app/modules/base"
	housedto "app/app/modules/house/dto"
	"app/internal/modules/log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type HouseController struct {
	houseSvc *HouseService
}

func newController(houseService *HouseService) *HouseController {
	return &HouseController{
		houseSvc: houseService,
	}
}

func (*HouseController) Get(ctx *gin.Context) {
	base.Success(ctx, "ok")
}

func (ctl *HouseController) CreateHouse(ctx *gin.Context) {
	req := housedto.HouseRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		log.Info("Error: %v", err)
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	imageMainFile, err := ctx.FormFile("image_main")
	if err != nil || imageMainFile == nil {
		log.Info("Error: %v", err)
		base.BadRequest(ctx, "Main image is required", nil)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		log.Info("Error: %v", err)
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	imageFiles := form.File["image"]
	if len(imageFiles) == 0 {
		base.BadRequest(ctx, "No image uploaded", nil)
		return
	}

	createdByID, exists := ctx.Get("userId")
	if !exists {
		base.BadRequest(ctx, "User not found", nil)
		return
	}

	createdByType, typeexist := ctx.Get("userType")
	if !typeexist {
		base.BadRequest(ctx, "User type not found", nil)
		return
	}

	houseResponse, err := ctl.houseSvc.CreateHouse(ctx, &req, imageMainFile, imageFiles, createdByID.(string), createdByType.(string))
	if err != nil {
		log.Info("Error: %v", err)
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponse)
}

func (ctl *HouseController) GetHouseByID(ctx *gin.Context) {
	houseID := ctx.Param("id")

	houseResponse, err := ctl.houseSvc.GetHouseByID(ctx, houseID)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponse)
}

func (ctl *HouseController) GetHousesByProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userId")
	if !exists {
		base.BadRequest(ctx, "User not found", nil)
		return
	}

	houseResponses, err := ctl.houseSvc.GetHousesByProfile(ctx, userID.(string))
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponses)
}

func (ctl *HouseController) GetAllHouses(ctx *gin.Context) {
	req := housedto.HouseGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	houseResponse, count, err := ctl.houseSvc.GetAllHouses(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, houseResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *HouseController) GetAllHousesForAdmin(ctx *gin.Context) {
	req := housedto.HouseGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	houseResponse, count, err := ctl.houseSvc.GetAllHousesForAdmin(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, houseResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *HouseController) GetHousesConfirmation(ctx *gin.Context) {
	req := housedto.HouseGetAllRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	log.Info("Request: %v", req)

	houseResponse, count, err := ctl.houseSvc.GetHousesConfirmation(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, houseResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *HouseController) UpdateHouse(ctx *gin.Context) {
	houseID := ctx.Param("id")
	req := housedto.HouseUpdateRequest{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}
	log.Info("Request: %v", req)

	form, err := ctx.MultipartForm()
	if err != nil {
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

	log.Info("ImageMain file: %v", ImageMainFile)
	log.Info("Image files: %v", imageFiles)

	err = ctl.houseSvc.UpdateHouse(ctx, houseID, &req, ImageMainFile, imageFiles)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, "success")
}

func (ctl *HouseController) UpdateRecommendHouse(ctx *gin.Context) {
	houseID := ctx.Param("id")
	req := housedto.HouseUpdateRecommendRequest{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	houseResponse, err := ctl.houseSvc.UpdateRecommendHouse(ctx, houseID, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponse)
}

func (ctl *HouseController) UpdateStatusHouse(ctx *gin.Context) {
	houseID := ctx.Param("id")
	req := housedto.HouseUpdateStatusRequest{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	houseResponse, err := ctl.houseSvc.UpdateStatusHouse(ctx, houseID, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponse)
}

func (ctl *HouseController) UpdateConfirmation(ctx *gin.Context) {
	houseID := ctx.Param("id")
	req := housedto.HouseUpdateConfirmationRequest{}

	if err := ctx.Bind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	houseResponse, err := ctl.houseSvc.UpdateConfirmation(ctx, houseID, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponse)
}

func (ctl *HouseController) DeleteHouse(ctx *gin.Context) {
	houseID := ctx.Param("id")

	err := ctl.houseSvc.DeleteHouse(ctx, houseID)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, nil)
}

func (ctl *HouseController) GetNearbyHouses(ctx *gin.Context) {

	houseResponses, err := ctl.houseSvc.GetNearbyHouses(ctx)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, houseResponses)
}

func (ctl *HouseController) GetHouseHistory(ctx *gin.Context) {
	req := housedto.HouseGetHistoryRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, "Invalid request", nil)
		return
	}

	houseResponse, count, err := ctl.houseSvc.GetHouseHistory(ctx, &req)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Paginate(ctx, houseResponse, &base.ResponsePaginate{
		From:  int64(req.From),
		Size:  int64(req.Size),
		Total: int64(count),
	})
}

func (ctl *HouseController) GetHouseCountByZone(ctx *gin.Context) {
	zoneCountResponses, err := ctl.houseSvc.GetHouseCountByZone(ctx)
	if err != nil {
		base.BadRequest(ctx, err.Error(), nil)
		return
	}

	base.Success(ctx, zoneCountResponses)
}

func (ctl *HouseController) GetPriceRange(ctx *gin.Context) {
	response, err := ctl.houseSvc.GetPriceRange(ctx)
	if err != nil {
		log.Info("Error: %v", err)
		base.BadRequest(ctx, "Failed to fetch price range", nil)
		return
	}

	base.Success(ctx, response)
}
