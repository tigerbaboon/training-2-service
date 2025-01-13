package activitylogs

import (
	activitylogsent "app/app/modules/activitylogs/ent"
	"context"

	"github.com/uptrace/bun"
)

type ActivitylogsService struct {
	db *bun.DB
}

func newService(db *bun.DB) *ActivitylogsService {
	return &ActivitylogsService{
		db: db,
	}
}

func (s *ActivitylogsService) CreateLogs(ctx context.Context, req activitylogsent.ActivityLogs) (*activitylogsent.ActivityLogs, error) {
	_, err := s.db.NewInsert().Model(&req).Exec(ctx)
	return &req, err
}
