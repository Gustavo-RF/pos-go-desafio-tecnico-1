package redisconfig

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Ctx         context.Context
	RedisClient redis.Client
}

func NewRedisClient(context context.Context, address, port, password string) *RedisConfig {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", address, port),
		Password: password,
		DB:       0,
	})

	return &RedisConfig{
		Ctx:         context,
		RedisClient: *rdb,
	}
}

func (r *RedisConfig) SetKey(key, value string, expireTime int) error {
	err := r.RedisClient.Set(r.Ctx, key, value, time.Duration(expireTime)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisConfig) GetKey(key string) (*string, error) {
	val2, err := r.RedisClient.Get(r.Ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
		return nil, fmt.Errorf("%s does not exist", key)
	} else if err != nil {
		return nil, err
	} else {
		fmt.Println("key2", val2)
		return &val2, nil
	}
}
