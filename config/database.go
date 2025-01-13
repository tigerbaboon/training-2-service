package config

import (
	dbdto "app/internal/modules/database/dto"
	rddto "app/internal/modules/redis/dto"
)

type Database struct {
	Sql   map[string]*dbdto.Option
	Redis map[string]*rddto.Option
}

var (
	database = Database{
		Sql: map[string]*dbdto.Option{"": {
			Host:     "127.0.0.1",
			Port:     5432,
			Database: "postgres",
			Username: "postgres",
			Password: "",
			TimeZone: "Asia/Bangkok",
		}},
		Redis: map[string]*rddto.Option{"": {
			Db:       0,
			Addr:     "",
			Username: "",
			Password: "",
		}},
	}
)
