package travelent

import "time"

type Travels struct {
	ID                string     `bun:"id,default:gen_random_uuid,pk"`
	TravelTitle       string     `bun:"travel_title"`
	TravelDetail      string     `bun:"travel_detail"`
	Status            bool       `bun:"status,default:true"`
	Address           string     `bun:"address"`
	LocationLatitute  float64    `bun:"location_latitute,notnull"`
	LocationLongitute float64    `bun:"location_longitute,notnull"`
	CreatedAt         int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt         int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedBy         string     `bun:"updated_by"`
	DeletedAt         *time.Time `bun:"deleted_at,soft_delete"`
}
