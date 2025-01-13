package log

import (
	logdto "app/app/modules/log/dto"
	logent "app/app/modules/log/ent"
	"app/internal/modules/log"
	"context"

	"github.com/uptrace/bun"
)

type LogService struct {
	db *bun.DB
}

func newService(db *bun.DB) *LogService {
	return &LogService{
		db: db,
	}
}

func (service *LogService) CreateLogs(ctx context.Context, req *logdto.LogDTORequest) (*logent.LogEntity, error) {
	logs := &logent.LogEntity{
		MenagerID:   req.MenagerID,
		ActionType:  req.ActionType,
		Description: req.Description,
		RecordID:    req.RecordID,
		TableName:   req.TableName,
		UpdatedBy:   req.UpdatedBy,
	}

	_, err := service.db.
		NewInsert().
		Model(logs).
		Exec(ctx)

	if err != nil {
		log.Info(err.Error())
		return nil, err
	}

	return logs, nil
}

func (svc *LogService) GetAllLog(ctx context.Context, limit int, offset int) ([]*logent.LogEntity, error) {
	var logs []*logent.LogEntity
	err := svc.db.NewSelect().
		Model(&logs).
		OrderExpr("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		log.Info("Error fetching all logs: %s", err.Error())
		return nil, err
	}
	return logs, nil
}

func (svc *LogService) DeleteLog(ctx context.Context, id string) error {
	_, err := svc.db.NewDelete().
		Model(&logent.LogEntity{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	return nil
}