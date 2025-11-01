package repositories

import (
	"context"
	"fmt"

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
			Addr: fmt.Sprintf("%s:%d", cfg.RedisCfg.Host, cfg.RedisCfg.Port),
			DB:   cfg.RedisCfg.DB,
		}),
	}
}

func (i IDGeneratorRepository) NextID() int {
	ctx := context.Background()

	result := i.conn.Incr(ctx, i.HashKey)
	return int(result.Val())
}
