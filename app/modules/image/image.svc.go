package image

import (
	imagedto "app/app/modules/image/dto"
	imageent "app/app/modules/image/ent"
	"context"
	"time"

	"github.com/uptrace/bun"
)

type ImageService struct {
	db *bun.DB
}

func newService(db *bun.DB) *ImageService {
	return &ImageService{
		db: db,
	}
}

func (service *ImageService) GetAllImages(ctx context.Context) ([]*imageent.Images, error) {
	var images []*imageent.Images

	err := service.db.
		NewSelect().
		Model(&images).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (service *ImageService) CreateImage(ctx context.Context, imgUrl, id string) error {
	image := &imageent.Images{
		StationID: id,
		ImageURL:  imgUrl,
	}

	_, err := service.db.
		NewInsert().
		Model(image).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *ImageService) CreateImagesWithType(ctx context.Context, imgUrl, imgType, id string) error {
	image := &imageent.Images{
		StationID: id,
		ImageURL:  imgUrl,
		Type:      imgType,
	}

	_, err := service.db.
		NewInsert().
		Model(image).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *ImageService) UpdateImage(ctx context.Context, id string, req *imagedto.ImageRequest) (*imageent.Images, error) {
	var image imageent.Images
	err := service.db.NewSelect().
		Model(&image).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if req.ImageURL != "" {
		image.ImageURL = req.ImageURL
	}

	image.UpdatedAt = time.Now().Unix()

	_, err = service.db.NewUpdate().
		Model(&image).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (svc *ImageService) GetImagesByID(ctx context.Context, sID string) ([]*imageent.Images, error) {
	var images []*imageent.Images

	err := svc.db.NewSelect().
		Model(&images).
		Where("s_id = ?", sID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return images, nil
}

func (svc *ImageService) DeleteImageByID(ctx context.Context, id string) error {
	_, err := svc.db.NewDelete().
		Model(&imageent.Images{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
