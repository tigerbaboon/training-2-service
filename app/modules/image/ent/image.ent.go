package imageent

import "time"

type Images struct {
	ID        string     `bun:"id,default:gen_random_uuid,pk"`
	StationID string     `bun:"s_id"`
	ImageURL  string     `bun:"image_url,notnull"`
	Type      string     `bun:"type,notnull"`
	CreatedAt int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedBy string     `bun:"updated_by"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}
