package house

import (
	"app/app/modules/amenity"
	"app/app/modules/contact"
	contactent "app/app/modules/contact/ent"
	housedto "app/app/modules/house/dto"
	houseent "app/app/modules/house/ent"
	"app/app/modules/image"
	"app/app/modules/manager"
	"app/app/modules/user"
	"app/app/modules/zone"
	helper "app/helper/googleStorage"
	"app/internal/modules/log"
	"context"
	"mime/multipart"

	"strings"

	"github.com/jinzhu/copier"
	"github.com/uptrace/bun"
)

type HouseService struct {
	db             *bun.DB
	ImageService   *image.ImageService
	ContactService *contact.ContactService
	AmenityService *amenity.AmenityService
	ZoneService    *zone.ZoneService
	UserService    *user.UserService
	ManagerService *manager.ManagerService
}

func newService(db *bun.DB, ImageService *image.ImageService, ContactService *contact.ContactService, AmenityService *amenity.AmenityService, ZoneService *zone.ZoneService, UserService *user.UserService, ManagerService *manager.ManagerService) *HouseService {
	return &HouseService{
		db:             db,
		ImageService:   ImageService,
		ContactService: ContactService,
		AmenityService: AmenityService,
		ZoneService:    ZoneService,
		UserService:    UserService,
		ManagerService: ManagerService,
	}
}

func (svc *HouseService) CreateHouse(ctx context.Context, req *housedto.HouseRequest, imageMainFile *multipart.FileHeader, imageFiles []*multipart.FileHeader, createdByID string, createdByType string) (*housedto.HouseDTOResponse, error) {
	Amenity_ID := []string{}
	if len(req.AmenityID) > 0 {
		Amenity_ID = strings.Split(string(req.AmenityID), ",")
	}
	house := &houseent.Houses{
		HouseName:         req.HouseName,
		HouseType:         req.HouseType,
		ZoneID:            req.ZoneID,
		SellType:          req.SellType,
		Size:              req.Size,
		AmenityID:         Amenity_ID,
		Floor:             req.Floor,
		Price:             req.Price,
		NumberOfRooms:     req.NumberOfRooms,
		NumberOfBathrooms: req.NumberOfBathrooms,
		WaterRate:         req.WaterRate,
		ElectricityRate:   req.ElectricityRate,
		Description:       req.Description,
		Address:           req.Address,
		LocationLatitute:  req.LocationLatitute,
		LocationLongitute: req.LocationLongitute,
		IsRecommend:       req.IsRecommend,
		CreatedBy:         createdByID,
		CreatedByType:     createdByType,
		Confirmation:      "pending",
	}

	_, err := svc.db.NewInsert().
		Model(house).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	contacts := &contactent.Contacts{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		LineID:      req.LineID,
		HouseID:     house.ID,
	}

	err = svc.ContactService.CreateContact(ctx, contacts, house.ID)
	if err != nil {
		return nil, err
	}

	imageMainURL := ""
	if imageMainFile != nil {
		imgURLPtr, err := helper.UploadFileGCS(imageMainFile)
		if err != nil {
			return nil, err
		}

		if imgURLPtr != nil {
			imageMainURL = *imgURLPtr
		}

		err = svc.ImageService.CreateImagesWithType(ctx, imageMainURL, "cover", house.ID)
		if err != nil {
			return nil, err
		}
	}

	imageURLs := []string{}
	for _, file := range imageFiles {
		imgURLPtr, err := helper.UploadFileGCS(file)
		if err != nil {
			return nil, err
		}

		var imgURL string
		if imgURLPtr != nil {
			imgURL = *imgURLPtr
		}

		imageURLs = append(imageURLs, imgURL)

		err = svc.ImageService.CreateImagesWithType(ctx, imgURL, "original", house.ID)
		if err != nil {
			return nil, err
		}
	}

	houseResponse := &housedto.HouseDTOResponse{}
	copier.Copy(houseResponse, house)

	houseResponse.ImagesMain = imageMainURL
	houseResponse.Images = imageURLs

	houseResponse.ContactInfo = housedto.ContactInfoResponse{
		FirstName:   contacts.FirstName,
		LastName:    contacts.LastName,
		PhoneNumber: contacts.PhoneNumber,
		LineID:      contacts.LineID,
	}

	houseResponse.AmenityID = Amenity_ID

	return houseResponse, nil
}

func (svc *HouseService) GetHouseByID(ctx context.Context, id string) (*housedto.HouseResponse, error) {
	var house housedto.HouseResponse

	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id", "h.house_name", "h.house_type", "h.sell_type", "h.size", "h.floor", "h.price", "h.number_of_rooms", "h.number_of_bathrooms", "h.water_rate", "h.electricity_rate", "h.description", "h.address", "h.location_latitute", "h.location_longitute", "h.is_recommend", "h.created_at", "h.confirmation", "h.amenity_id").
		ColumnExpr("z.id as zone__id, z.zone_name as zone__zone_name").
		ColumnExpr("c.first_name as contact_info__first_name, c.last_name as contact_info__last_name, c.phone_number as contact_info__phone_number, c.line_id as contact_info__line_id").
		ColumnExpr("i.id as images_main__id, i.image_url as images_main__url, i.type as images_main__type").
		Join("LEFT JOIN zone_entities as z ON h.zone_id = z.id").
		Join("LEFT JOIN contacts as c ON h.id = c.house_id").
		Join("LEFT JOIN images as i ON h.id = i.s_id AND i.type = ? AND i.deleted_at IS NULL", "cover").
		Join("LEFT JOIN amenity_entities as a ON h.amenity_id = a.id").
		Where("h.id = ?", id)

	err := query.Scan(ctx, &house)
	if err != nil {
		return nil, err
	}

	var images []housedto.ImageResponse
	err = svc.db.NewSelect().
		TableExpr("images as i").
		Column("i.id").
		ColumnExpr("i.image_url as url").
		Column("i.type").
		Where("i.s_id = ? AND i.type = ?", id, "original").
		Where("i.deleted_at IS NULL").
		Scan(ctx, &images)
	if err != nil {
		return nil, err
	}
	house.Images = images

	var amenities []housedto.AmenityResponse
	err = svc.db.NewSelect().
		TableExpr("amenity_entities as a").
		Column("a.id", "a.icons").
		ColumnExpr("a.amenity_name as name").
		Where("a.id IN (?)", bun.In(house.AmenityID)).
		Scan(ctx, &amenities)
	if err != nil {
		return nil, err
	}
	house.Amenity = amenities
	house.AmenityID = nil

	return &house, nil
}

func (svc *HouseService) GetHousesByProfile(ctx context.Context, userid string) ([]housedto.HouseResponseForGetByProfile, error) {
	var houses []housedto.HouseResponseForGetByProfile

	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id", "h.house_name", "h.address", "h.size", "h.floor", "h.price", "h.number_of_rooms", "h.number_of_bathrooms", "h.status", "h.created_at").
		ColumnExpr("c.first_name as contact_info__first_name").
		ColumnExpr("c.last_name as contact_info__last_name").
		ColumnExpr("c.phone_number as contact_info__phone_number").
		ColumnExpr("c.line_id as contact_info__line_id").
		ColumnExpr("i.id as images_main__id").
		ColumnExpr("i.image_url as images_main__url").
		Join("left join contacts as c ON h.id = c.house_id").
		Join("left join images as i ON h.id = i.s_id AND i.type = ?", "cover").
		Where("h.created_by = ?", userid)

	err := query.Scan(ctx, &houses)
	if err != nil {
		return nil, err
	}

	return houses, nil
}

func (svc *HouseService) GetAllHouses(ctx context.Context, req *housedto.HouseGetAllRequest) ([]*housedto.HouseResponseForGetALL, int, error) {
	var houses []*housedto.HouseResponseForGetALL

	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id", "h.house_name", "h.house_type", "h.sell_type", "h.size", "h.floor", "h.price", "h.number_of_rooms", "h.number_of_bathrooms", "h.address", "h.is_recommend").
		ColumnExpr("z.id as zone__id").
		ColumnExpr("z.zone_name as zone__zone_name").
		ColumnExpr("i.id as images_main__id").
		ColumnExpr("i.image_url as images_main__url").
		ColumnExpr("i.type as images_main__type").
		ColumnExpr("c.first_name as contact_info__first_name").
		ColumnExpr("c.last_name as contact_info__last_name").
		ColumnExpr("c.phone_number as contact_info__phone_number").
		ColumnExpr("c.line_id as contact_info__line_id").
		Join("left join zone_entities as z on h.zone_id = z.id").
		Join("left join images as i on h.id = i.s_id AND i.type = ? AND i.deleted_at IS NULL", "cover").
		Join("left join contacts as c on h.id = c.house_id").
		Where("h.confirmation = ? AND h.deleted_at IS NULL", "approved")

	if req.SearchByName != "" {
		query.Where("house_name LIKE ?", "%"+req.SearchByName+"%")
	}
	if req.SearchByZone != "" {
		query.Where("zone_id = ?", req.SearchByZone)
	}
	if req.PriceStart != 0 && req.PriceEnd != 0 {
		query.Where("price BETWEEN ? AND ?", req.PriceStart, req.PriceEnd)
	}

	query.Order("h.is_recommend DESC").OrderExpr("RANDOM()")

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &houses)
	if err != nil {
		return nil, 0, err
	}

	return houses, count, nil
}

func (svc *HouseService) GetAllHousesForAdmin(ctx context.Context, req *housedto.HouseGetAllRequest) ([]*housedto.HouseResponseForAdmin, int, error) {
	var houses []*housedto.HouseResponseForAdmin

	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id", "h.house_name", "h.house_type", "h.is_recommend").
		Where("h.deleted_at IS NULL AND h.confirmation = ?", "approved")

	if req.SearchByName != "" {
		query.Where("h.house_name LIKE ?", "%"+req.SearchByName+"%")
	}

	query.Order("h.is_recommend DESC")
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &houses)
	if err != nil {
		return nil, 0, err
	}

	return houses, count, nil
}

func (svc *HouseService) GetHousesConfirmation(ctx context.Context, req *housedto.HouseGetAllRequest) ([]*housedto.HouseResponseForConfirmation, int, error) {
	var houses []*housedto.HouseResponseForConfirmation

	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id", "h.house_name", "h.house_type", "h.confirmation").
		Where("h.confirmation = ?", "pending")

	if req.SearchByName != "" {
		query.Where("h.house_name LIKE ?", "%"+req.SearchByName+"%")
	}

	query.Order("h.created_at DESC")
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &houses)
	if err != nil {
		return nil, 0, err
	}

	return houses, count, nil
}

func (svc *HouseService) UpdateHouse(ctx context.Context, id string, req *housedto.HouseUpdateRequest, imageMainFile *multipart.FileHeader, imageFiles []*multipart.FileHeader) error {
	house := &houseent.Houses{}
	err := svc.db.NewSelect().
		Model(house).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return err
	}

	if req.HouseName != "" {
		house.HouseName = req.HouseName
	}
	if req.HouseType != "" {
		house.HouseType = req.HouseType
	}
	if req.ZoneID != "" {
		house.ZoneID = req.ZoneID
	}
	if req.SellType != "" {
		house.SellType = req.SellType
	}
	if req.Size != 0 {
		house.Size = req.Size
	}
	if req.Floor != 0 {
		house.Floor = req.Floor
	}
	if req.Price != 0 {
		house.Price = req.Price
	}
	if req.NumberOfRooms != 0 {
		house.NumberOfRooms = req.NumberOfRooms
	}
	if req.NumberOfBathrooms != 0 {
		house.NumberOfBathrooms = req.NumberOfBathrooms
	}
	if req.WaterRate != 0 {
		house.WaterRate = req.WaterRate
	}
	if req.ElectricityRate != 0 {
		house.ElectricityRate = req.ElectricityRate
	}
	if req.Description != "" {
		house.Description = req.Description
	}
	if req.Address != "" {
		house.Address = req.Address
	}
	if req.LocationLatitute != 0 {
		house.LocationLatitute = req.LocationLatitute
	}
	if req.LocationLongitute != 0 {
		house.LocationLongitute = req.LocationLongitute
	}

	house.Confirmation = "pending"

	_, err = svc.db.NewUpdate().
		Model(house).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	contact := &contactent.Contacts{}
	err = svc.db.NewSelect().
		Model(contact).
		Where("house_id = ?", id).
		Scan(ctx)
	if err == nil {
		if req.FirstName != "" {
			contact.FirstName = req.FirstName
		}
		if req.LastName != "" {
			contact.LastName = req.LastName
		}
		if req.PhoneNumber != "" {
			contact.PhoneNumber = req.PhoneNumber
		}
		if req.LineID != "" {
			contact.LineID = req.LineID
		}

		_, err = svc.db.NewUpdate().
			Model(contact).
			Where("house_id = ?", id).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	oldImages, err := svc.ImageService.GetImagesByID(ctx, house.ID)
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
			err = svc.ImageService.CreateImagesWithType(ctx, *imgURLPtr, "cover", house.ID)
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
			err = svc.ImageService.CreateImagesWithType(ctx, *imgURLPtr, "original", house.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (svc *HouseService) UpdateRecommendHouse(ctx context.Context, id string, req *housedto.HouseUpdateRecommendRequest) (*housedto.HouseUpdateRecommendResponse, error) {
	house := &houseent.Houses{}
	err := svc.db.NewSelect().
		Model(house).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	house.IsRecommend = req.IsRecommend

	_, err = svc.db.NewUpdate().
		Model(house).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	houseResponse := &housedto.HouseUpdateRecommendResponse{}
	copier.Copy(houseResponse, house)

	return houseResponse, nil
}

func (svc *HouseService) UpdateStatusHouse(ctx context.Context, id string, req *housedto.HouseUpdateStatusRequest) (*housedto.HouseUpdateStatusResponse, error) {
	house := &houseent.Houses{}
	err := svc.db.NewSelect().
		Model(house).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	house.Status = req.Status

	_, err = svc.db.NewUpdate().
		Model(house).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	houseResponse := &housedto.HouseUpdateStatusResponse{}
	copier.Copy(houseResponse, house)

	return houseResponse, nil
}

func (svc *HouseService) UpdateConfirmation(ctx context.Context, id string, req *housedto.HouseUpdateConfirmationRequest) (*housedto.HouseUpdateConfirmationResponse, error) {
	house := &houseent.Houses{}
	err := svc.db.NewSelect().
		Model(house).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	house.Confirmation = req.Confirmation

	_, err = svc.db.NewUpdate().
		Model(house).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	houseResponse := &housedto.HouseUpdateConfirmationResponse{}
	copier.Copy(houseResponse, house)

	return houseResponse, nil
}

func (svc *HouseService) DeleteHouse(ctx context.Context, id string) error {
	_, err := svc.db.NewDelete().
		TableExpr("contacts").
		Where("house_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = svc.db.NewDelete().
		TableExpr("images").
		Where("s_id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = svc.db.NewDelete().
		TableExpr("houses").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (svc *HouseService) GetNearbyHouses(ctx context.Context) ([]*housedto.HouseNearByResponse, error) {
	var houses []*housedto.HouseNearByResponse
	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id",
			"h.house_name",
			"h.house_type",
			"h.sell_type",
			"h.size",
			"h.price",
			"h.number_of_rooms",
			"h.number_of_bathrooms",
			"h.address",
			"h.location_latitute",
			"h.location_longitute")

	err := query.Scan(ctx, &houses)
	if err != nil {
		return nil, err
	}

	return houses, nil
}

func (svc *HouseService) GetHouseHistory(ctx context.Context, req *housedto.HouseGetHistoryRequest) ([]*housedto.HouseResponseHistory, int, error) {
	var houses []*housedto.HouseResponseHistory

	query := svc.db.NewSelect().
		TableExpr("houses as h").
		Column("h.id", "h.house_name", "h.house_type", "h.sell_type", "h.created_at").
		ColumnExpr(`
			CASE 
				WHEN h.created_by_type = 'user' THEN CONCAT(u.firstname, ' ', u.lastname)
				WHEN h.created_by_type = 'manager' THEN m.manager_name
				ELSE 'Unknown'
			END AS created_by
		`).
		Join("LEFT JOIN user_entities as u ON h.created_by = u.id AND h.created_by_type = 'user'").
		Join("LEFT JOIN managers as m ON h.created_by = m.id AND h.created_by_type = 'manager'").
		Where("h.deleted_at IS NULL")

	if req.SearchByName != "" {
		query.Where("h.house_name LIKE ?", "%"+req.SearchByName+"%")
	}

	query.Order("h.created_at DESC")

	count, err := query.Count(ctx)
	if err != nil {
		log.Info("error : %v", err)
		return nil, 0, err
	}

	if req.From == 1 {
		req.From = 0
	}

	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &houses)
	if err != nil {
		log.Info("error : %v", err)
		return nil, 0, err
	}

	return houses, count, nil
}

func (svc *HouseService) GetHouseCountByZone(ctx context.Context) ([]housedto.ZoneCountResponse, error) {
	var zoneCountResponses []housedto.ZoneCountResponse

	err := svc.db.NewSelect().
		TableExpr("houses as h").
		ColumnExpr("z.zone_name, COUNT(h.id) as count").
		Join("LEFT JOIN zone_entities as z ON h.zone_id = z.id").
		Group("z.zone_name").
		Order("z.zone_name ASC").
		Where("h.confirmation = ?", "approved").
		Where("h.deleted_at IS NULL").
		Scan(ctx, &zoneCountResponses)
	if err != nil {
		return nil, err
	}

	return zoneCountResponses, nil
}

func (svc *HouseService) GetPriceRange(ctx context.Context) (*housedto.PriceRangeResponse, error) {
	response := &housedto.PriceRangeResponse{}

	err := svc.db.NewSelect().
		Model((*houseent.Houses)(nil)).
		ColumnExpr("MIN(price) AS price_min, MAX(price) AS price_max").
		Scan(ctx, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
// 	const earthRadiusKm = 6371
// 	latRad1 := lat1 * math.Pi / 180       // แปลง lat1 จากองศาเป็นเรเดียน
// 	latRad2 := lng2 * math.Pi / 180       // แปลง lat2 จากองศาเป็นเรเดียน
// 	dLat := (lat2 - lat1) * math.Pi / 180 // หาระยะทางระหว่าง lat2 และ lat1
// 	dLng := (lng2 - lng1) * math.Pi / 180 // หาระยะทางระหว่าง lng2 และ lng1
// 	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
// 		math.Cos(latRad1)*math.Cos(latRad2)*
// 			math.Sin(dLng/2)*math.Sin(dLng/2) // ใช้สูตร Haversine คำนวณค่า a (ค่า a ใช้สำหรับคำนวณระยะทางบนผิวโลกตามสูตร Haversine)
// 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)) // คำนวณค่า c จากค่า a และค่า 1-a (ค่า c คือมุม (หน่วยเรเดียน) ระหว่างจุดสองจุดบนผิวโลก)
// 	return earthRadiusKm * c // คำนวณระยะทางระหว่างจุดสองจุดบนผิวโลก
// }

// func (svc *HouseService) GetNearbyHouses(ctx context.Context, req *housedto.NearbyHouses) ([]*housedto.HouseResponse, error) {
// 	const radius = 4.0
// 	latRadius := radius / 111.0                                        // 1 องศาเท่ากับ 111 กิโลเมตร
// 	lngRadius := radius / (111.0 * math.Cos(req.Latitude*math.Pi/180)) // 1 องศาเท่ากับ 111 กิโลเมตร และ 1 องศาเท่ากับ 111 กิโลเมตร แต่ต้องคูณด้วยค่า Cos ของ latitude
// 	var houses []*houseent.Houses
// 	err := svc.db.NewSelect().
// 		Model(&houses).
// 		Where("location_latitute BETWEEN ? AND ?", req.Latitude-latRadius, req.Latitude+latRadius).
// 		Where("location_longitute BETWEEN ? AND ?", req.Longitude-lngRadius, req.Longitude+lngRadius).
// 		Scan(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var nearbyHouses []*housedto.HouseResponse
// 	for _, house := range houses {
// 		distance := calculateDistance(req.Latitude, req.Longitude, house.LocationLatitute, house.LocationLongitute)
// 		if distance <= radius {
// 			houseResponse := &housedto.HouseResponse{}
// 			copier.Copy(houseResponse, house)
// 			contact := &contactent.Contacts{}
// 			err = svc.db.NewSelect().
// 				Model(contact).
// 				Where("house_id = ?", house.ID).
// 				Scan(ctx)
// 			if err != nil {
// 				return nil, err
// 			} else {
// 				houseResponse.ContactInfo = housedto.ContactInfoResponse{
// 					FirstName:   contact.FirstName,
// 					LastName:    contact.LastName,
// 					PhoneNumber: contact.PhoneNumber,
// 					LineID:      contact.LineID,
// 				}
// 			}
// 			var amenities []amenityent.AmenityEntity
// 			err = svc.db.NewSelect().
// 				Model(&amenities).
// 				Where("id IN (?)", bun.In(house.AmenityID)).
// 				Scan(ctx)
// 			if err != nil {
// 				return nil, err
// 			} else {
// 				houseResponse.Amenity = []housedto.AmenityResponse{}
// 				for _, amenity := range amenities {
// 					houseResponse.Amenity = append(houseResponse.Amenity, housedto.AmenityResponse{
// 						ID:   amenity.ID,
// 						Name: amenity.AmenityName,
// 					})
// 				}
// 			}
// 			zoneResponse, err := svc.ZoneService.GetZoneByID(ctx, house.ZoneID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			houseResponse.Zone = housedto.ZoneResponse{
// 				ID:       zoneResponse.ID,
// 				ZoneName: zoneResponse.ZoneName,
// 			}
// 			images, err := svc.ImageService.GetImagesByID(ctx, house.ID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			imageMain := housedto.ImageResponse{}
// 			imageOriginal := []housedto.ImageResponse{}
// 			for _, img := range images {
// 				imageResponse := housedto.ImageResponse{
// 					ID:   img.ID,
// 					URL:  img.ImageURL,
// 					Type: img.Type,
// 				}
// 				if img.Type == "cover" {
// 					imageMain = imageResponse
// 				} else {
// 					imageOriginal = append(imageOriginal, imageResponse)
// 				}
// 			}
// 			houseResponse.ImagesMain = imageMain
// 			houseResponse.Images = imageOriginal
// 			nearbyHouses = append(nearbyHouses, houseResponse)
// 		}
// 	}
// 	return nearbyHouses, nil
// }

// func (svc *HouseService) GetAllHouses(ctx context.Context, req *housedto.HouseGetAllRequest) ([]*housedto.HouseResponse, int, error) {
// 	var houses []*houseent.Houses
// 	query := svc.db.NewSelect().
// 		Model(&houses)
// 	if req.SearchByName != "" {
// 		query.Where("house_name LIKE ?", "%"+req.SearchByName+"%")
// 	}
// 	if req.SearchByZone != "" {
// 		query.Where("zone_id = ?", req.SearchByZone)
// 	}
// 	if req.PriceStart != 0 && req.PriceEnd != 0 {
// 		query.Where("price BETWEEN ? AND ?", req.PriceStart, req.PriceEnd)
// 	}
// 	query.Order("is_recommend DESC").OrderExpr("RANDOM()")
// 	err := query.Scan(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	count, err := query.Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	if req.From == 1 {
// 		req.From = 0
// 	}
// 	err = query.Offset(req.From).Limit(req.Size).Scan(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	var houseResponses []*housedto.HouseResponse
// 	for _, house := range houses {
// 		houseResponse := &housedto.HouseResponse{}
// 		copier.Copy(houseResponse, house)
// 		contact := &contactent.Contacts{}
// 		err = svc.db.NewSelect().Model(contact).Where("house_id = ?", house.ID).Scan(ctx)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		houseResponse.ContactInfo = housedto.ContactInfoResponse{
// 			FirstName:   contact.FirstName,
// 			LastName:    contact.LastName,
// 			PhoneNumber: contact.PhoneNumber,
// 			LineID:      contact.LineID,
// 		}
// 		var amenities []amenityent.AmenityEntity
// 		err = svc.db.NewSelect().
// 			Model(&amenities).
// 			Where("id IN (?)", bun.In(house.AmenityID)).
// 			Scan(ctx)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		log.Info("amenities : %v", house.AmenityID)
// 		houseResponse.Amenity = []housedto.AmenityResponse{}
// 		for _, amenity := range amenities {
// 			houseResponse.Amenity = append(houseResponse.Amenity, housedto.AmenityResponse{
// 				ID:   amenity.ID,
// 				Name: amenity.AmenityName,
// 			})
// 		}
// 		zoneResponse, err := svc.ZoneService.GetZoneByID(ctx, house.ZoneID)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		houseResponse.Zone = housedto.ZoneResponse{
// 			ID:       zoneResponse.ID,
// 			ZoneName: zoneResponse.ZoneName,
// 		}
// 		images, err := svc.ImageService.GetImagesByID(ctx, house.ID)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		imageMain := housedto.ImageResponse{}
// 		imageOriginal := []housedto.ImageResponse{}
// 		for _, img := range images {
// 			imageResponse := housedto.ImageResponse{
// 				ID:   img.ID,
// 				URL:  img.ImageURL,
// 				Type: img.Type,
// 			}
// 			if img.Type == "cover" {
// 				imageMain = imageResponse
// 			} else {
// 				imageOriginal = append(imageOriginal, imageResponse)
// 			}
// 		}
// 		houseResponse.ImagesMain = imageMain
// 		houseResponse.Images = imageOriginal
// 		houseResponses = append(houseResponses, houseResponse)
// 	}
// 	return houseResponses, count, nil
// }

// func (svc *HouseService) GetAllHouses(ctx context.Context, req *housedto.HouseGetAllRequest) ([]*housedto.HouseResponse, int, error) {
// 	var houses []*housedto.HouseResponse
// 	query := svc.db.NewSelect().
// 		TableExpr("houses as h").
// 		Column("h.id", "h.house_name", "h.house_type", "h.sell_type", "h.size", "h.floor", "h.price", "h.number_of_rooms", "h.number_of_bathrooms", "h.water_rate", "h.electricity_rate", "h.description", "h.address", "h.location_latitute", "h.location_longitute", "h.is_recommend", "h.created_at").
// 		ColumnExpr("z.id as zone__id").
// 		ColumnExpr("z.zone_name as zone__zone_name").
// 		ColumnExpr("i.id as images_main__id").
// 		ColumnExpr("i.image_url as images_main__url").
// 		ColumnExpr("i.type as images_main__type").
// 		ColumnExpr("c.first_name as contact_info__first_name").
// 		ColumnExpr("c.last_name as contact_info__last_name").
// 		ColumnExpr("c.phone_number as contact_info__phone_number").
// 		ColumnExpr("c.line_id as contact_info__line_id").
// 		Join("left join zone_entities as z on h.zone_id = z.id").
// 		Join("left join images as i on h.id = i.s_id").
// 		Join("left join contacts as c on h.id = c.house_id").
// 		Where("i.type = ?", "cover")
// 	if req.SearchByName != "" {
// 		query.Where("house_name LIKE ?", "%"+req.SearchByName+"%")
// 	}
// 	if req.SearchByZone != "" {
// 		query.Where("zone_id = ?", req.SearchByZone)
// 	}
// 	if req.PriceStart != 0 && req.PriceEnd != 0 {
// 		query.Where("price BETWEEN ? AND ?", req.PriceStart, req.PriceEnd)
// 	}
// 	query.Order("is_recommend DESC").OrderExpr("RANDOM()")
// 	count, err := query.Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	if req.From == 1 {
// 		req.From = 0
// 	}
// 	err = query.Offset(req.From).Limit(req.Size).Scan(ctx, &houses)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	var img []housedto.ImageResponse
// 	for _, house := range houses {
// 		err = svc.db.NewSelect().
// 			TableExpr("images as i").
// 			Column("i.id", "i.type").
// 			ColumnExpr("i.image_url AS url").
// 			Where("i.s_id = ?", house.ID).
// 			Where("i.type = ?", "original").
// 			Scan(ctx, &img)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		house.Images = img
// 		var amtID string
// 		err = svc.db.NewSelect().
// 			TableExpr("houses as h").
// 			Column("h.amenity_id").
// 			Scan(ctx, &amtID)
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 		amtID = strings.ReplaceAll(amtID, `"`, "") // ลบเครื่องหมายคำพูดทั้งหมด
// 		amtID = strings.Trim(amtID, "[]")          // ลบ [] ด้านนอก
// 		ids := strings.Split(amtID, ",")           // แยกค่าโดยใช้คอมม่า
// 		var amenities []housedto.AmenityResponse
// 		err = svc.db.NewSelect().
// 			TableExpr("amenity_entities").
// 			Column("id").
// 			ColumnExpr("amenity_name as name").
// 			ColumnExpr("icons").
// 			Where("id IN (?)", bun.In(ids)).
// 			Scan(ctx, &amenities)
// 		if err != nil {
// 			log.Info("error1 : %v", err)
// 			return nil, 0, err
// 		}
// 		house.Amenity = amenities
// 	}
// 	return houses, count, nil
// }
