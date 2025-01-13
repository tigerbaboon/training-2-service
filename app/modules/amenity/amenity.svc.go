package amenity

import (
	amenitydto "app/app/modules/amenity/dto"
	amenityent "app/app/modules/amenity/ent"
	"app/internal/modules/log"
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/uptrace/bun"
)

type AmenityService struct {
	db *bun.DB
}

func newService(db *bun.DB) *AmenityService {
	return &AmenityService{
		db: db,
	}
}

func (service *AmenityService) CreateAmenity(ctx context.Context, req *amenitydto.AmenityRequest) error {
	amenity := &amenityent.AmenityEntity{
		AmenityName: req.AmenityName,
		Icons:       req.Icons,
	}

	_, err := service.db.
		NewInsert().
		Model(amenity).
		Exec(ctx)

	return err
}

func (service *AmenityService) GetAmenityByID(ctx context.Context, id string) (*amenityent.AmenityEntity, error) {
	amenity := &amenityent.AmenityEntity{}
	err := service.db.NewSelect().
		Model(amenity).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return amenity, nil
}

func (service *AmenityService) GetAllAmenity(ctx context.Context, req *amenitydto.AmenityGetAllRequest) ([]*amenitydto.AmenityDTOResponse, int, error) {
	var amenity []*amenityent.AmenityEntity

	query := service.db.NewSelect().
		Model(&amenity)
	if req.Search != "" {
		query.Where("amenity_name LIKE ?", "%"+req.Search+"%")
	}
	query.Order("created_at ASC")
	count, err := query.Count(ctx)

	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	var amenityResponses []*amenitydto.AmenityDTOResponse
	copier.CopyWithOption(&amenityResponses, &amenity, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return amenityResponses, count, nil
}

func (service *AmenityService) UpdateAmenity(ctx context.Context, id string, req *amenitydto.AmenityRequest) (*amenitydto.AmenityDTOResponse, error) {
	amenity := &amenityent.AmenityEntity{}
	err := service.db.NewSelect().
		Model(amenity).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	amenity.AmenityName = req.AmenityName
	amenity.Icons = req.Icons

	_, err = service.db.NewUpdate().
		Model(amenity).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	AmenityResponse := &amenitydto.AmenityDTOResponse{}
	copier.Copy(AmenityResponse, amenity)

	return AmenityResponse, nil
}

func (service *AmenityService) DeleteAmenity(ctx context.Context, id string) error {
	exists, err := service.db.NewSelect().
		TableExpr("houses").
		Where("? = ANY (amenity_id)", id).
		Where("deleted_at IS NULL").
		Exists(ctx)
	if err != nil {
		log.Info(err.Error())
		return err
	}

	if exists {
		return errors.New("amenity is being used")
	}

	_, err = service.db.NewDelete().
		Model(&amenityent.AmenityEntity{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	return nil
}
