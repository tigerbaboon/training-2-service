package redis

import (
	"app/internal/modules/config"
	dto "app/internal/modules/redis/dto"
	"context"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	rMap map[string]*JSONClient
	mut  sync.RWMutex
}

var (
	defaultRDConfig = redis.Options{
		MaxRetries: 1,
	}
)

func newService(confService *config.ConfigService) *RedisService {
	conf := confService.App()
	service := &RedisService{
		rMap: make(map[string]*JSONClient),
	}
	if err := service.Register(context.Background(), conf.AppEnv, conf.Database.Redis); err != nil {
		panic(err)
	}
	return service
}

func (rds *RedisService) DB(name ...string) *JSONClient {
	rds.mut.RLock()
	defer rds.mut.RUnlock()
	if rds.rMap == nil {
		panic("redis not initialized")
	}
	if len(name) == 0 {
		return rds.rMap[""]
	}

	db, ok := rds.rMap[name[0]]
	if !ok {
		panic("redis not initialized")
	}
	return db
}

func (rds *RedisService) Register(ctx context.Context, appName string, opts map[string]*dto.Option) error {
	rds.mut.Lock()
	defer rds.mut.Unlock()
	for key, opt := range opts {
		optDef := withRDDefaultConf(appName, opt)
		rd := redis.NewClient(optDef)
		rdJson := newJSONClient(rd)
		if err := rdJson.Ping(ctx).Err(); err != nil {
			return fmt.Errorf("redis ping error: %w", err)
		}
		rds.rMap[key] = rdJson
	}
	return nil
}

func (rds *RedisService) Close(ctx context.Context, name string) error {
	rds.mut.Lock()
	defer rds.mut.Unlock()
	if rds.rMap == nil {
		return nil
	}
	rd, ok := rds.rMap[name]
	if !ok {
		return nil
	}
	if err := rd.Close(); err != nil {
		return err
	}
	delete(rds.rMap, name)
	return nil

}

func (rds *RedisService) close(ctx context.Context) error {
	rds.mut.Lock()
	defer rds.mut.Unlock()
	for _, db := range rds.rMap {
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func withRDDefaultConf(appName string, opt *dto.Option) *redis.Options {
	rOpt := defaultRDConfig
	rOpt.ClientName = appName
	copier.CopyWithOption(&rOpt, opt, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return &rOpt
}
