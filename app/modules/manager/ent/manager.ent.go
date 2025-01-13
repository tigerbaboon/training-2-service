package managerent

import "time"

type Managers struct {
	ID          string     `bun:"id,default:gen_random_uuid,pk"`
	Username    string     `bun:"username,unique"`
	Password    string     `bun:"password,notnull"`
	ManagerName string     `bun:"manager_name,notnull"`
	CreatedAt   int64      `bun:"created_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedAt   int64      `bun:"updated_at,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	UpdatedBy   string     `bun:"updated_by"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
