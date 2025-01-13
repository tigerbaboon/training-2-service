package zone

import (
	"app/app/modules/image"
	imageent "app/app/modules/image/ent"
	zonedto "app/app/modules/zone/dto"
	zoneent "app/app/modules/zone/ent"
	helper "app/helper/googleStorage"
	"app/internal/modules/log"
	"context"
	"errors"
	"mime/multipart"

	"github.com/uptrace/bun"
)

type ZoneService struct {
	db           *bun.DB
	ImageService *image.ImageService
}

func newService(db *bun.DB, ImageService *image.ImageService) *ZoneService {
	return &ZoneService{
		db:           db,
		ImageService: ImageService,
	}
}

func (service *ZoneService) CreateZone(ctx context.Context, req *zonedto.ZoneRequest, files []*multipart.FileHeader) error {
	zone := &zoneent.ZoneEntity{
		ZoneName: req.ZoneName,
		Lat:      req.Lat,
		Long:     req.Long,
	}

	_, err := service.db.NewInsert().Model(zone).Exec(ctx)
	if err != nil {
		log.Info("Error creating zone: %s", err.Error())
		return err
	}

	for _, file := range files {
		imgURLPtr, err := helper.UploadFileGCS(file)
		if err != nil {
			return err
		}

		var imgURL string
		if imgURLPtr != nil {
			imgURL = *imgURLPtr
		} else {
			imgURL = ""
		}

		err = service.ImageService.CreateImagesWithType(ctx, imgURL, "zone", zone.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *ZoneService) GetZoneByID(ctx context.Context, id string) (*zonedto.ZoneResponse, error) {
	var zone zonedto.ZoneResponse

	query := service.db.NewSelect().
		TableExpr("zone_entities as z").
		Column("z.id", "z.zone_name", "z.lat", "z.long").
		ColumnExpr("i.id as images__id").
		ColumnExpr("i.image_url as images__url").
		ColumnExpr("i.type as images__type").
		Join("left join images as i ON z.id = i.s_id AND i.type = ? AND i.deleted_at IS NULL", "zone").
		Where("z.id = ?", id)

	err := query.Scan(ctx, &zone)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (service *ZoneService) GetAllZone(ctx context.Context, req *zonedto.ZoneGetAllRequest) ([]*zonedto.ZoneResponses, int, error) {
	var zones []*zonedto.ZoneResponses

	query := service.db.NewSelect().
		TableExpr("zone_entities as z").
		Column("z.id", "z.zone_name", "z.lat", "z.long").
		ColumnExpr("i.id as images__id").
		ColumnExpr("i.image_url as images__url").
		ColumnExpr("i.type as images__type").
		Join("LEFT JOIN images as i ON z.id = i.s_id AND i.type = ?", "zone").
		Where("i.deleted_at IS NULL").
		Where("z.deleted_at IS NULL")

	if req.Search != "" {
		query.Where("z.zone_name LIKE ?", "%"+req.Search+"%")
	}

	query.Order("z.created_at ASC")
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &zones)
	if err != nil {
		return nil, 0, err
	}

	return zones, count, nil
}

func (service *ZoneService) UpdateZone(ctx context.Context, id string, req *zonedto.ZoneUpdateRequest, files []*multipart.FileHeader) error {
	zones := &zoneent.ZoneEntity{}
	err := service.db.NewSelect().
		Model(zones).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	zones.ZoneName = req.ZoneName
	zones.Lat = req.Lat
	zones.Long = req.Long

	_, err = service.db.NewUpdate().
		Model(zones).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	var url *string

	if len(files) > 0 {
		imgURL, err := helper.UploadFileGCS(files[0])
		if err != nil {
			return err
		}
		url = imgURL
	}

	for i := range req.RemainingImageIDs {
		img := &imageent.Images{}
		err := service.db.NewSelect().
			Model(img).
			Where("id = ?", req.RemainingImageIDs[i]).
			Scan(ctx)
		if err != nil {
			return err
		}

		if url != nil {
			img.ImageURL = *url
		}

		_, err = service.db.NewUpdate().
			Model(img).
			Where("id = ?", req.RemainingImageIDs[i]).
			Exec(ctx)
		if err != nil {
			return err
		}

	}

	return nil
}

func (service *ZoneService) DeleteZone(ctx context.Context, id string) error {
	exists, err := service.db.NewSelect().
		TableExpr("houses").
		Where("zone_id = ?", id).
		Where("deleted_at IS NULL").
		Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		log.Info("Zone is still being used by a house")
		return errors.New("zone is still being used by a house")
	}

	_, err = service.db.NewDelete().
		Model(&zoneent.ZoneEntity{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
