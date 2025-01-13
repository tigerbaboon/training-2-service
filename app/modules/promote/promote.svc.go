package promote

import (
	"app/app/modules/image"
	imageent "app/app/modules/image/ent"
	promotedto "app/app/modules/promote/dto"
	promoteent "app/app/modules/promote/ent"
	helper "app/helper/googleStorage"
	"context"
	"mime/multipart"

	"github.com/jinzhu/copier"
	"github.com/uptrace/bun"
)

type PromoteService struct {
	db           *bun.DB
	ImageService *image.ImageService
}

func newService(db *bun.DB, ImageService *image.ImageService) *PromoteService {
	return &PromoteService{
		db:           db,
		ImageService: ImageService,
	}
}

func (service *PromoteService) CreatePromote(ctx context.Context, req *promotedto.PromoteDTORequest, files []*multipart.FileHeader) (*promotedto.PromoteDTOResponse, error) {
	promote := &promoteent.Promotes{
		PromoteName: req.PromoteName,
		PromoteType: req.PromoteType,
		Link:        req.Link,
		Status:      req.Status,
	}

	_, err := service.db.NewInsert().
		Model(promote).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	imageURLs := []string{}

	for _, file := range files {
		imgURLPtr, err := helper.UploadFileGCS(file)
		if err != nil {
			return nil, err
		}

		var imgURL string
		if imgURLPtr != nil {
			imgURL = *imgURLPtr
		} else {
			imgURL = ""
		}

		imageURLs = append(imageURLs, imgURL)

		err = service.ImageService.CreateImagesWithType(ctx, imgURL, "promote", promote.ID)
		if err != nil {
			return nil, err
		}
	}

	promoteResponse := &promotedto.PromoteDTOResponse{}
	copier.Copy(promoteResponse, promote)
	promoteResponse.Images = imageURLs

	return promoteResponse, nil
}

func (service *PromoteService) GetPromoteByID(ctx context.Context, id string) (*promotedto.PromoteResponse, error) {
	promote := &promoteent.Promotes{}
	err := service.db.NewSelect().
		Model(promote).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	promoteResponse := &promotedto.PromoteResponse{}
	copier.Copy(promoteResponse, promote)

	images, err := service.ImageService.GetImagesByID(ctx, promote.ID)
	if err != nil {
		return nil, err
	}

	imageResponses := []promotedto.ImageResponse{}
	for _, img := range images {
		imageResponses = append(imageResponses, promotedto.ImageResponse{
			ID:  img.ID,
			URL: img.ImageURL,
		})
	}

	promoteResponse.Images = imageResponses

	return promoteResponse, nil
}

func (service *PromoteService) GetAllPromotes(ctx context.Context, req *promotedto.PromoteGetAllRequest) ([]*promotedto.PromoteResponseAll, int, error) {
	var promotes []*promotedto.PromoteResponseAll

	query := service.db.NewSelect().
		TableExpr("promotes as p").
		Column("p.id", "p.promote_name", "p.promote_type", "p.link", "p.status").
		ColumnExpr("i.id as image__id, i.image_url as image__url, i.type as image__type").
		Join("LEFT JOIN images as i ON p.id = i.s_id AND i.type = ? AND i.deleted_at IS NULL", "promote").
		Where("p.deleted_at IS NULL")

	if req.Search != "" {
		query.Where("p.promote_name LIKE ?", "%"+req.Search+"%")
	}
	if req.Type != "" {
		query.Where("p.promote_type = ?", req.Type)
	}

	query.OrderExpr("RANDOM()")
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &promotes)
	if err != nil {
		return nil, 0, err
	}

	return promotes, count, nil
}

func (service *PromoteService) UpdatePromote(ctx context.Context, id string, req *promotedto.PromoteUpdateRequest, imageFile *multipart.FileHeader) error {
	promote := &promoteent.Promotes{}
	err := service.db.NewSelect().
		Model(promote).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	if req.PromoteName != "" {
		promote.PromoteName = req.PromoteName
	}
	if req.PromoteType != "" {
		promote.PromoteType = req.PromoteType
	}
	if req.Link != "" {
		promote.Link = req.Link
	}

	_, err = service.db.NewUpdate().
		Model(promote).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	var url *string

	if imageFile != nil {
		imgURL, err := helper.UploadFileGCS(imageFile)
		if err != nil {
			return err
		}
		url = imgURL
	}

	img := &imageent.Images{}
	err = service.db.NewSelect().
		Model(img).
		Where("s_id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	if url != nil {
		img.ImageURL = *url
	}

	_, err = service.db.NewUpdate().
		Model(img).
		Where("s_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *PromoteService) UpdatePromoteStatus(ctx context.Context, id string, req *promotedto.PromoteUpdateStatus) (*promotedto.PromoteDTOResponse, error) {
	promote := &promoteent.Promotes{}
	err := service.db.NewSelect().
		Model(promote).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	promote.Status = req.Status

	_, err = service.db.NewUpdate().
		Model(promote).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	images, err := service.ImageService.GetImagesByID(ctx, promote.ID)
	if err != nil {
		return nil, err
	}

	imgURLs := []string{}
	for _, image := range images {
		imgURLs = append(imgURLs, image.ImageURL)
	}

	promoteResponse := &promotedto.PromoteDTOResponse{}
	copier.Copy(promoteResponse, promote)
	promoteResponse.Images = imgURLs

	return promoteResponse, nil
}

func (service *PromoteService) DeletePromote(ctx context.Context, id string) error {
	promote := &promoteent.Promotes{}
	err := service.db.NewSelect().
		Model(promote).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	images := []*imageent.Images{}
	err = service.db.NewSelect().
		Model(&images).
		Where("s_id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	for _, img := range images {
		err = service.ImageService.DeleteImageByID(ctx, img.ID)
		if err != nil {
			return err
		}
	}

	_, err = service.db.NewDelete().
		Model(promote).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
