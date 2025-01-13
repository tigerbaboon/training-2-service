package houseent

import "time"

type Houses struct {
	ID                string     `bun:"id,default:gen_random_uuid,pk"`
	HouseName         string     `bun:"house_name,notnull"`
	HouseType         string     `bun:"house_type,notnull"`
	ZoneID            string     `bun:"zone_id,notnull"`
	SellType          string     `bun:"sell_type,notnull"`
	AmenityID         []string   `json:"amenity_id"`
	Size              float64    `bun:"size,notnull"`
	Floor             float64    `bun:"floor,notnull"`
	Price             float64    `bun:"price,notnull"`
	NumberOfRooms     int64      `bun:"number_of_rooms,notnull"`
	NumberOfBathrooms int64      `bun:"number_of_bathrooms,notnull"`
	WaterRate         float64    `bun:"water_rate,notnull"`
	ElectricityRate   float64    `bun:"electricity_rate,notnull"`
	Description       string     `bun:"description,notnull"`
	Address           string     `bun:"address,notnull"`
	LocationLatitute  float64    `bun:"location_latitute,notnull"`
	LocationLongitute float64    `bun:"location_longitute,notnull"`
	IsRecommend       bool       `bun:"is_recommend"`
	CreatedAt         int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt         int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedBy         string     `bun:"updated_by"`
	DeletedAt         *time.Time `bun:"deleted_at,soft_delete"`
	CreatedBy         string     `bun:"created_by"`
	CreatedByType     string     `bun:"created_by_type"`
	Status            bool       `bun:"status,default:true"`
	Confirmation      string     `bun:"confirmation"`
}
