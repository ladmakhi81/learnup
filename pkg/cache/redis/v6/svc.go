package redisv6

import (
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/pkg/cache"
)

type RedisClientSvc struct {
	redis *redis.Client
}

func NewRedisClientSvc(redis *redis.Client) *RedisClientSvc {
	return &RedisClientSvc{
		redis: redis,
	}
}

func (svc *RedisClientSvc) SetHashVal(key string, id string, val any) error {
	err := svc.redis.HSet(key, id, val)
	if err != nil {
		return cache.NewCacheError(
			"Error: happen in set value",
			"RedisClientSvc.SetHashVal",
		)
	}
	return nil
}

func (svc *RedisClientSvc) GetHashVal(key, id string) (string, error) {
	val, err := svc.redis.HGet(key, id).Result()
	if err != nil {
		if err == redis.Nil {
			return "", cache.NewCacheError(
				"Error: cache key not found",
				"RedisClientSvc.GetHashVal",
			)
		}
		return "", cache.NewCacheError(
			"Error: happen in get value",
			"RedisClientSvc.GetHashVal",
		)
	}
	return val, nil
}

func (svc *RedisClientSvc) SetVal(key string, val any) error {
	err := svc.redis.Set(key, val, 0).Err()
	if err != nil {
		return cache.NewCacheError(
			"Error: happen in set value",
			"RedisClientSvc.SetVal",
		)
	}
	return nil
}
