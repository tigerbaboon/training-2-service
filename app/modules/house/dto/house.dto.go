package housedto

type HouseRequest struct {
	HouseName         string  `form:"house_name" binding:"required"`
	HouseType         string  `form:"house_type" binding:"required"`
	ZoneID            string  `form:"zone_id" binding:"required"`
	SellType          string  `form:"sell_type" binding:"required"`
	AmenityID         string  `form:"amenity_id" binding:"required"`
	Size              float64 `form:"size" binding:"required"`
	Floor             float64 `form:"floor" binding:"required"`
	Price             float64 `form:"price" binding:"required"`
	NumberOfRooms     int64   `form:"number_of_rooms" binding:"required"`
	NumberOfBathrooms int64   `form:"number_of_bathrooms" binding:"required"`
	WaterRate         float64 `form:"water_rate" binding:"required"`
	ElectricityRate   float64 `form:"electricity_rate" binding:"required"`
	Description       string  `form:"description"`
	Address           string  `form:"address" binding:"required"`
	LocationLatitute  float64 `form:"location_latitute" binding:"required"`
	LocationLongitute float64 `form:"location_longitute" binding:"required"`
	IsRecommend       bool    `form:"is_recommend"`
	FirstName         string  `form:"first_name" binding:"required"`
	LastName          string  `form:"last_name" binding:"required"`
	PhoneNumber       string  `form:"phone_number"`
	LineID            string  `form:"line_id"`
}

type ContactInfoResponse struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	LineID      string `json:"line_id"`
}

type HouseResponse struct {
	ID                string              `json:"id"`
	HouseName         string              `json:"house_name"`
	HouseType         string              `json:"house_type"`
	Zone              ZoneResponse        `json:"zone"`
	SellType          string              `json:"sell_type"`
	Size              float64             `json:"size"`
	Amenity           []AmenityResponse   `json:"amenity"`
	AmenityID         []string            `json:"amenity_id"`
	Floor             float64             `json:"floor"`
	Price             float64             `json:"price"`
	NumberOfRooms     int64               `json:"number_of_rooms"`
	NumberOfBathrooms int64               `json:"number_of_bathrooms"`
	WaterRate         float64             `json:"water_rate"`
	ElectricityRate   float64             `json:"electricity_rate"`
	Description       string              `json:"description"`
	Address           string              `json:"address"`
	LocationLatitute  float64             `json:"location_latitute"`
	LocationLongitute float64             `json:"location_longitute"`
	IsRecommend       bool                `json:"is_recommend"`
	ImagesMain        ImageResponse       `json:"images_main"`
	Images            []ImageResponse     `json:"images"`
	ContactInfo       ContactInfoResponse `json:"contact_info"`
	CreatedAt         int64               `json:"created_at"`
	Confirmation      string              `json:"confirmation"`
}

type HouseResponseForAdmin struct {
	ID          string `json:"id"`
	HouseName   string `json:"house_name"`
	HouseType   string `json:"house_type"`
	IsRecommend bool   `json:"is_recommend"`
}

type HouseResponseForConfirmation struct {
	ID           string `json:"id"`
	HouseName    string `json:"house_name"`
	HouseType    string `json:"house_type"`
	Confirmation string `json:"confirmation"`
}

type HouseResponseForGetALL struct {
	ID                string              `json:"id"`
	HouseName         string              `json:"house_name"`
	HouseType         string              `json:"house_type"`
	Zone              ZoneResponse        `json:"zone"`
	SellType          string              `json:"sell_type"`
	Size              float64             `json:"size"`
	Floor             float64             `json:"floor"`
	Price             float64             `json:"price"`
	NumberOfRooms     int64               `json:"number_of_rooms"`
	NumberOfBathrooms int64               `json:"number_of_bathrooms"`
	Address           string              `json:"address"`
	IsRecommend       bool                `json:"is_recommend"`
	ImagesMain        ImageResponse       `json:"images_main"`
	ContactInfo       ContactInfoResponse `json:"contact_info"`
}

type HouseResponseForGetByProfile struct {
	ID                string              `json:"id"`
	HouseName         string              `json:"house_name"`
	Address           string              `json:"address"`
	Size              float64             `json:"size"`
	Floor             float64             `json:"floor"`
	Price             float64             `json:"price"`
	NumberOfRooms     int64               `json:"number_of_rooms"`
	NumberOfBathrooms int64               `json:"number_of_bathrooms"`
	ContactInfo       ContactInfoResponse `json:"contact_info"`
	ImagesMain        ImageResponse       `json:"images_main"`
	Status            bool                `json:"status"`
	CreatedAt         int64               `json:"created_at"`
}

type ZoneResponse struct {
	ID       string `json:"id"`
	ZoneName string `json:"zone_name"`
}
type AmenityResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icons string `json:"icons"`
}

type HouseDTOResponse struct {
	ID                string              `json:"id"`
	HouseName         string              `json:"house_name"`
	HouseType         string              `json:"house_type"`
	ZoneID            string              `json:"zone_id"`
	SellType          string              `json:"sell_type"`
	Size              float64             `json:"size"`
	AmenityID         []string            `json:"amenity_id"`
	Floor             float64             `json:"floor"`
	Price             float64             `json:"price"`
	NumberOfRooms     int64               `json:"number_of_rooms"`
	NumberOfBathrooms int64               `json:"number_of_bathrooms"`
	WaterRate         float64             `json:"water_rate"`
	ElectricityRate   float64             `json:"electricity_rate"`
	Description       string              `json:"description"`
	Address           string              `json:"address"`
	LocationLatitute  float64             `json:"location_latitute"`
	LocationLongitute float64             `json:"location_longitute"`
	IsRecommend       bool                `json:"is_recommend"`
	ImagesMain        string              `json:"images_main"`
	Images            []string            `json:"images"`
	ContactInfo       ContactInfoResponse `json:"contact_info"`
}

type HouseUpdateRequest struct {
	HouseName            string   `form:"house_name"`
	HouseType            string   `form:"house_type"`
	ZoneID               string   `form:"zone_id"`
	SellType             string   `form:"sell_type"`
	Size                 float64  `form:"size"`
	AmenityID            string   `form:"amenity_id"`
	Floor                float64  `form:"floor"`
	Price                float64  `form:"price"`
	NumberOfRooms        int64    `form:"number_of_rooms"`
	NumberOfBathrooms    int64    `form:"number_of_bathrooms"`
	WaterRate            float64  `form:"water_rate"`
	ElectricityRate      float64  `form:"electricity_rate"`
	Description          string   `form:"description"`
	Address              string   `form:"address"`
	LocationLatitute     float64  `form:"location_latitute"`
	LocationLongitute    float64  `form:"location_longitute"`
	FirstName            string   `form:"first_name"`
	LastName             string   `form:"last_name"`
	PhoneNumber          string   `form:"phone_number"`
	LineID               string   `form:"line_id"`
	RemainingImageIDs    []string `form:"remainingImageIDs"`
	RemainingImageMainID string   `form:"remainingImageMainID"`
}

type HouseResponseHistory struct {
	ID        string `json:"id"`
	HouseName string `json:"house_name"`
	HouseType string `json:"house_type"`
	SellType  string `json:"sell_type"`
	CreatedAt int64  `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type HouseGetAllRequest struct {
	From         int    `form:"from"`
	Size         int    `form:"size"`
	SearchByName string `form:"search_by_name"`
	SearchByZone string `form:"search_by_zone"`
	PriceStart   int    `form:"price_start"`
	PriceEnd     int    `form:"price_end"`
}

type HouseGetHistoryRequest struct {
	From         int    `form:"from"`
	Size         int    `form:"size"`
	SearchByName string `form:"search_by_name"`
}

type ImageResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

type HouseUpdateRecommendRequest struct {
	IsRecommend bool `json:"is_recommend"`
}

type HouseUpdateRecommendResponse struct {
	ID          string `json:"id"`
	HouseName   string `json:"house_name"`
	IsRecommend bool   `json:"is_recommend"`
}

type HouseUpdateStatusRequest struct {
	Status bool `json:"status"`
}

type HouseUpdateStatusResponse struct {
	ID        string `json:"id"`
	HouseName string `json:"house_name"`
	Status    bool   `json:"status"`
}

type HouseUpdateConfirmationRequest struct {
	Confirmation string `json:"confirmation"`
}

type HouseUpdateConfirmationResponse struct {
	ID           string `json:"id"`
	HouseName    string `json:"house_name"`
	Confirmation string `json:"confirmation"`
}

type ZoneCount struct {
	ZoneID string `bun:"zone_id"`
	Count  int    `bun:"count"`
}

type ZoneCountResponse struct {
	ZoneName string `json:"zone_name"`
	Count    int    `json:"count"`
}

type PriceRangeResponse struct {
	PriceMin float64 `json:"price_min"`
	PriceMax float64 `json:"price_max"`
}

type HouseNearByResponse struct {
	ID                string  `json:"id"`
	HouseName         string  `json:"house_name"`
	HouseType         string  `json:"house_type"`
	SellType          string  `json:"sell_type"`
	Size              float64 `json:"size"`
	Price             float64 `json:"price"`
	NumberOfRooms     int64   `json:"number_of_rooms"`
	NumberOfBathrooms int64   `json:"number_of_bathrooms"`
	Address           string  `json:"address"`
	LocationLatitute  float64 `json:"location_latitute"`
	LocationLongitute float64 `json:"location_longitute"`
}
