package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	pointsAllKey = "points:all"
)

type PointsCache struct {
	rdb *redis.Client
}

func NewPointsCache(rdb *redis.Client) *PointsCache {
	return &PointsCache{rdb: rdb}
}

func (c *PointsCache) GetAll(ctx context.Context) (string, error) {
	return c.rdb.Get(ctx, pointsAllKey).Result()
}

func (c *PointsCache) SetAll(ctx context.Context, data []byte) error {
	return c.rdb.Set(ctx, pointsAllKey, data, 5*time.Minute).Err()
}

func (c *PointsCache) InvalidateAll(ctx context.Context) error {
	return c.rdb.Del(ctx, pointsAllKey).Err()
}
