package amenityent

import "time"

type AmenityEntity struct {
	ID          string     `bun:"id,default:gen_random_uuid,pk"`
	AmenityName string     `bun:"amenity_name,notnull"`
	Icons       string     `bun:"icons`
	CreatedAt   int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt   int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
