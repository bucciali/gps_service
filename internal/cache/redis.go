package cache

import "github.com/redis/go-redis/v9"

func NewRedisClient(add string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: add,
	})
}
