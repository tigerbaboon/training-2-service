package zoneent

import "time"

type ZoneEntity struct {
	ID        string     `bun:"id,default:gen_random_uuid,pk"`
	ZoneName  string     `bun:"zone_name"`
	Lat       float64    `bun:"lat,notnull"`
	Long      float64    `bun:"long,notnull"`
	CreatedAt int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedBy string     `bun:"updated_by"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}
