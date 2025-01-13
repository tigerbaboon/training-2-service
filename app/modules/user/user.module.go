package user

import (
	"github.com/uptrace/bun"
)

type UserModule struct {
	Ctl     *UserCTL
	Service *UserService
}

func New(db *bun.DB) *UserModule {
	service := newService(db)
	return &UserModule{
		Ctl:     newCTL(service),
		Service: service,
	}
}
