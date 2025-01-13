package activitylogsent

import "github.com/uptrace/bun"

type ActivityLogs struct {
	bun.BaseModel `bun:"activitylogs"`

	ID            string      `bun:"default:gen_random_uuid(),pk"`
	Section       string      `bun:"section"`
	EventType     string      `bun:"eventtype"`
	StatusCode    string      `bun:"statuscode"`
	Detail        string      `bun:"detail"`
	Request       interface{} `json:"request" bun:"type:jsonb" form:"request"`
	Responses     interface{} `json:"responses" bun:"type:jsonb" form:"responses"`
	IpAddress     string      `bun:"ipaddress"`
	UserAgent     string      `bun:"useragent"`
	CreatedBy     string      `bun:"createdby"`
	CreatedAt     int64       `bun:"createdat,nullzero,default:EXTRACT(EPOCH FROM NOW())"`
	CreatedByType string      `bun:"createdbytype"`
}
