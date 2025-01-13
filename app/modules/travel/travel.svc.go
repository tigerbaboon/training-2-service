package travel

import (
	"app/app/modules/image"
	traveldto "app/app/modules/travel/dto"
	travelent "app/app/modules/travel/ent"
	helper "app/helper/googleStorage"
	storageHelper "app/helper/googleStorage"
	"context"
	"mime/multipart"
	"strings"

	"github.com/uptrace/bun"
)

type TravelService struct {
	db           *bun.DB
	ImageService *image.ImageService
}

func newService(db *bun.DB, ImageService *image.ImageService) *TravelService {
	return &TravelService{
		db:           db,
		ImageService: ImageService,
	}
}

func (svc *TravelService) CreateTravel(ctx context.Context, req *traveldto.TravelRequest, imageMainFile *multipart.FileHeader, imageFiles []*multipart.FileHeader) error {
	travel := &travelent.Travels{
		TravelTitle:       req.TravelTitle,
		TravelDetail:      req.TravelDetail,
		LocationLatitute:  req.LocationLatitute,
		LocationLongitute: req.LocationLongitute,
		Address:           req.Address,
	}

	_, err := svc.db.NewInsert().
		Model(travel).
		Exec(ctx)
	if err != nil {
		return err
	}

	if imageMainFile != nil {
		imgURLPtr, err := storageHelper.UploadFileGCS(imageMainFile)
		if err != nil {
			return err
		}

		if imgURLPtr != nil {
			err = svc.ImageService.CreateImagesWithType(ctx, *imgURLPtr, "cover", travel.ID)
			if err != nil {
				return err
			}
		}
	}

	for _, file := range imageFiles {
		imgURLPtr, err := storageHelper.UploadFileGCS(file)
		if err != nil {
			return err
		}

		if imgURLPtr != nil {
			err = svc.ImageService.CreateImagesWithType(ctx, *imgURLPtr, "original", travel.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (svc *TravelService) GetTravelByID(ctx context.Context, id string) (*traveldto.TravelResponseByID, error) {
	var travel traveldto.TravelResponseByID

	query := svc.db.NewSelect().
		TableExpr("travels as t").
		Column("t.id", "t.travel_title", "t.travel_detail", "t.location_latitute", "t.location_longitute", "t.address").
		ColumnExpr("i.id as image_main__id, i.image_url as image_main__url, i.type as image_main__type").
		Join("LEFT JOIN images as i ON t.id = i.s_id AND i.type = ?", "cover").
		Where("t.id = ?", id)

	err := query.Scan(ctx, &travel)
	if err != nil {
		return nil, err
	}

	var images []traveldto.ImageResponse
	err = svc.db.NewSelect().
		TableExpr("images as i").
		Column("i.id").
		ColumnExpr("i.image_url as url").
		Column("i.type").
		Where("i.s_id = ? AND i.type = ? AND deleted_at IS NULL", id, "original").
		Scan(ctx, &images)
	if err != nil {
		return nil, err
	}
	travel.Images = images

	return &travel, nil
}

func (svc *TravelService) GetAllTravels(ctx context.Context, req *traveldto.TravelGetAllRequest) ([]*traveldto.TravelResponseByGetAll, int, error) {
	var travels []*traveldto.TravelResponseByGetAll

	query := svc.db.NewSelect().
		TableExpr("travels as t").
		Column("t.id", "t.travel_title", "t.status").
		ColumnExpr("i.id as image_main__id").
		ColumnExpr("i.image_url as image_main__url").
		ColumnExpr("i.type as image_main__type").
		Join("left join images as i ON t.id = i.s_id AND i.type = ?", "cover").
		Where("t.deleted_at IS NULL")

	if req.Search != "" {
		query.Where("t.travel_title LIKE ?", "%"+req.Search+"%")
	}

	query.OrderExpr("RANDOM()")
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &travels)
	if err != nil {
		return nil, 0, err
	}

	return travels, count, nil
}

func (svc *TravelService) GetAllTravelForAdmin(ctx context.Context, req *traveldto.TravelGetAllRequest) ([]*traveldto.TravelResponseByID, int, error) {
	var travels []*traveldto.TravelResponseByID

	query := svc.db.NewSelect().
		TableExpr("travels as t").
		Column("t.id", "t.travel_title", "t.travel_detail", "t.status", "t.address", "t.location_latitute", "t.location_longitute").
		ColumnExpr("i.id as image_main__id, i.image_url as image_main__url, i.type as image_main__type").
		Join("LEFT JOIN images as i ON t.id = i.s_id AND i.type = ? AND i.deleted_at IS NULL", "cover").
		Where("t.deleted_at IS NULL")

	if req.Search != "" {
		query.Where("t.travel_title LIKE ?", "%"+req.Search+"%")
	}

	query.Order("t.created_at DESC")

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &travels)
	if err != nil {
		return nil, 0, err
	}

	for _, travel := range travels {
		var images []traveldto.ImageResponse
		err := svc.db.NewSelect().
			TableExpr("images as i").
			Column("i.id").
			ColumnExpr("i.image_url as url").
			Column("i.type").
			Where("i.s_id = ? AND i.type = ? AND deleted_at IS NULL", travel.ID, "original").
			Scan(ctx, &images)
		if err != nil {
			return nil, 0, err
		}
		travel.Images = images
	}

	return travels, count, nil
}

func (svc *TravelService) UpdateTravel(ctx context.Context, id string, req *traveldto.TravelUpdateRequest, imageMainFile *multipart.FileHeader, imageFiles []*multipart.FileHeader) error {
	travel := &travelent.Travels{}
	err := svc.db.NewSelect().
		Model(travel).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	if req.TravelTitle != "" {
		travel.TravelTitle = req.TravelTitle
	}
	if req.TravelDetail != "" {
		travel.TravelDetail = req.TravelDetail
	}
	if req.Address != "" {
		travel.Address = req.Address
	}
	if req.LocationLatitute != 0 {
		travel.LocationLatitute = req.LocationLatitute
	}
	if req.LocationLongitute != 0 {
		travel.LocationLongitute = req.LocationLongitute
	}

	_, err = svc.db.NewUpdate().
		Model(travel).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	oldImages, err := svc.ImageService.GetImagesByID(ctx, travel.ID)
	if err != nil {
		return err
	}

	remainingImageIDs := strings.Split(req.RemainingImageIDs[0], ",")
	remainingImageMap := make(map[string]bool)
	for _, imgID := range remainingImageIDs {
		remainingImageMap[imgID] = true
	}

	for _, img := range oldImages {
		if img.ID == req.RemainingImageMainID {
			continue
		}
		if !remainingImageMap[img.ID] {
			err = svc.ImageService.DeleteImageByID(ctx, img.ID)
			if err != nil {
				return err
			}
		}
	}

	if imageMainFile != nil {
		imgURLPtr, err := helper.UploadFileGCS(imageMainFile)
		if err != nil {
			return err
		}
		if imgURLPtr != nil {
			err = svc.ImageService.CreateImagesWithType(ctx, *imgURLPtr, "cover", travel.ID)
			if err != nil {
				return err
			}
		}
	}

	for _, file := range imageFiles {
		imgURLPtr, err := helper.UploadFileGCS(file)
		if err != nil {
			return err
		}
		if imgURLPtr != nil {
			err = svc.ImageService.CreateImagesWithType(ctx, *imgURLPtr, "original", travel.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (svc *TravelService) UpdateTravelStatus(ctx context.Context, id string, req *traveldto.TravelUpdateStatus) (*traveldto.TravelStatusResponse, error) {
	_, err := svc.db.NewUpdate().
		TableExpr("travels").
		Set("status = ?", req.Status).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	var travelResponse traveldto.TravelStatusResponse
	err = svc.db.NewSelect().
		TableExpr("travels").
		Column("id", "status").
		Where("id = ?", id).
		Scan(ctx, &travelResponse)
	if err != nil {
		return nil, err
	}

	return &travelResponse, nil
}

func (svc *TravelService) DeleteTravel(ctx context.Context, id string) error {
	_, err := svc.db.NewDelete().
		TableExpr("images").
		Where("s_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = svc.db.NewDelete().
		TableExpr("travels").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
