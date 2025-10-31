package repositories

import (
	"context"

	"github.com/coci/cutme/pkg/config"
	"github.com/redis/go-redis/v9"
)

type IDGeneratorRepository struct {
	conn    *redis.Client
	HashKey string
}

func NewIDGeneratorRepository(cfg *config.Config) *IDGeneratorRepository {
	return &IDGeneratorRepository{
		HashKey: cfg.RedisCfg.RedisHashKey,
		conn: redis.NewClient(&redis.Options{
			Addr:     cfg.RedisCfg.Host,
			Password: cfg.RedisCfg.Password,
			DB:       cfg.RedisCfg.DB,
		}),
	}
}

func (i IDGeneratorRepository) NextID() int {
	ctx := context.Background()

	result := i.conn.Incr(ctx, i.HashKey)
	return int(result.Val())
}
