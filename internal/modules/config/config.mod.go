package config

type ConfigModule struct {
	Svc *ConfigService
}

func New() *ConfigModule {
	svc := newService()
	return &ConfigModule{
		Svc: svc,
	}
}
