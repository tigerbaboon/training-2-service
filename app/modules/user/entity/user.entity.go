package entity

import (
	"time"
)

type UserEntity struct {
	ID        string     `bun:"id,default:gen_random_uuid,pk"`
	Username  string     `bun:"username,unique"`
	Password  string     `bun:"password,notnull"`
	Email     string     `bun:"email,notnull"`
	FirstName string     `bun:"firstname,notnull"`
	LastName  string     `bun:"lastname,notnull"`
	CreatedAt int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}
