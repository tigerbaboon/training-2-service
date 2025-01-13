package promoteent

import "time"

type Promotes struct {
	ID          string     `bun:"id,default:gen_random_uuid,pk"`
	PromoteName string     `bun:"promote_name,notnull"`
	PromoteType string     `bun:"promote_type,notnull"`
	Link        string     `bun:"link"`
	Status      bool       `bun:"status,default:true"`
	CreatedAt   int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt   int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedBy   string     `bun:"updated_by"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
