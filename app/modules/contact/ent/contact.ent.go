package contactent

import "time"

type Contacts struct {
	ID          string     `bun:"id,default:gen_random_uuid,pk"`
	HouseID     string     `bun:"house_id, notnull"`
	FirstName   string     `bun:"first_name,notnull"`
	LastName    string     `bun:"last_name,notnull"`
	PhoneNumber string     `bun:"phone_number"`
	LineID      string     `bun:"line_id"`
	CreatedAt   int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt   int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
