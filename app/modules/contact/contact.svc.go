package contact

import (
	contactent "app/app/modules/contact/ent"
	"context"

	"github.com/uptrace/bun"
)

type ContactService struct {
	db *bun.DB
}

func newService(db *bun.DB) *ContactService {
	return &ContactService{
		db: db,
	}
}

func (svc *ContactService) CreateContact(ctx context.Context, contactx *contactent.Contacts, houseID string) error {
	contactx.HouseID = houseID

	_, err := svc.db.NewInsert().Model(contactx).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ContactService) DeleteContact(ctx context.Context, id string) error {
	_, err := svc.db.NewDelete().
		Model(&contactent.Contacts{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
