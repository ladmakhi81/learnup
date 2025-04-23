package redisv6

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ladmakhi81/learnup/pkg/dtos"
)

type RedisClientSvc struct {
	redis *redis.Client
}

func setupRedisClient(config *dtos.EnvConfig) *redis.Client {
	host := config.Redis.Host
	port := config.Redis.Port
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})
}

func NewRedisClientSvc(config *dtos.EnvConfig) *RedisClientSvc {
	client := setupRedisClient(config)
	return &RedisClientSvc{
		redis: client,
	}
}

func (svc RedisClientSvc) SetHashVal(key string, id string, val any) error {
	err := svc.redis.HSet(key, id, val)
	if err != nil {
		return dtos.NewCacheError(
			"Error: happen in set value",
			"RedisClientSvc.SetHashVal",
		)
	}
	return nil
}

func (svc RedisClientSvc) GetHashVal(key, id string) (string, error) {
	val, err := svc.redis.HGet(key, id).Result()
	if err != nil {
		if err == redis.Nil {
			return "", dtos.NewCacheError(
				"Error: cache key not found",
				"RedisClientSvc.GetHashVal",
			)
		}
		return "", dtos.NewCacheError(
			"Error: happen in get value",
			"RedisClientSvc.GetHashVal",
		)
	}
	return val, nil
}

func (svc RedisClientSvc) SetVal(key string, val any) error {
	err := svc.redis.Set(key, val, 0).Err()
	if err != nil {
		return dtos.NewCacheError(
			"Error: happen in set value",
			"RedisClientSvc.SetVal",
		)
	}
	return nil
}

func (svc RedisClientSvc) GetVal(key string) (string, error) {
	val, err := svc.redis.Get(key).Result()
	if err != nil {
		return "", dtos.NewCacheError("Error: happen in get value", "RedisClientSvc.GetVal")
	}
	return val, nil
}
