package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	kiosksAllKey = "terminal_points:all"
)

type KioskCache struct {
	rdb *redis.Client
}

func NewKioskCache(rdb *redis.Client) *KioskCache {
	return &KioskCache{rdb: rdb}
}

func (c *KioskCache) GetAll(ctx context.Context) (string, error) {
	return c.rdb.Get(ctx, kiosksAllKey).Result()
}

func (c *KioskCache) SetAll(ctx context.Context, data []byte) error {
	return c.rdb.Set(ctx, kiosksAllKey, data, 5*time.Minute).Err()
}

func (c *KioskCache) InvalidateAll(ctx context.Context) error {
	return c.rdb.Del(ctx, kiosksAllKey).Err()
}
