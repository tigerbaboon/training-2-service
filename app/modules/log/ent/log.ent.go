package logent

import "time"

type LogEntity struct {
	ID          string     `bun:"id,default:gen_random_uuid,pk"`
	MenagerID   string     `bun:"menager_id"`
	ActionType  string     `bun:"actiontype"`
	Description string     `bun:"description"`
	RecordID    string     `bun:"record_id"`
	TableName   string     `bun:"table_name"`
	CreatedAt   int64      `bun:"created_at"`
	UpdatedAt   int64      `bun:"updated_at"`
	UpdatedBy   string     `bun:"updated_by"`
	DeletedAt   *time.Time `bun:"deleted_at,soft_delete"`
}
