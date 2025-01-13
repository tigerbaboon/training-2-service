package googlestorage

import "app/internal/modules/config"

type GoogleStorageModule struct {
	Svc *GoogleStorageService
}

func New(conf *config.ConfigService) *GoogleStorageModule {
	svc := newService(conf.App())
	return &GoogleStorageModule{
		Svc: svc,
	}
}
