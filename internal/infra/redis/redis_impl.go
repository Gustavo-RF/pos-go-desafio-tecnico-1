package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Ctx                  context.Context
	RedisClient          redis.Client
	RateLimiter          string
	RequestsPerSecond    int
	BlockedTimeInSeconds int
}

func NewRedisClient(ctx context.Context, address, port, password, rateLimiter string, request, blocked int) *RedisConfig {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", address, port),
		Password: password,
		DB:       0,
	})

	if rateLimiter == "" {
		rateLimiter = "IP"
	}

	if rateLimiter != "IP" && rateLimiter != "TOKEN" {
		panic(errors.New("rate limiter should be IP or TOKEN"))
	}

	return &RedisConfig{
		Ctx:                  ctx,
		RedisClient:          *rdb,
		RateLimiter:          rateLimiter,
		RequestsPerSecond:    request,
		BlockedTimeInSeconds: blocked,
	}
}

func (r *RedisConfig) SetKey(key, value string, expireTime int) error {
	err := r.RedisClient.Set(r.Ctx, key, value, time.Second*time.Duration(expireTime)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisConfig) GetKey(key string) (*string, error) {
	val2, err := r.RedisClient.Get(r.Ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%s does not exist", key)
	} else if err != nil {
		return nil, err
	} else {
		return &val2, nil
	}
}

func (r *RedisConfig) KeyExists(key string) bool {
	val, err := r.GetKey(key)
	if err != nil {
		return false
	}

	return val != nil
}

func (r *RedisConfig) IncKey(key string) error {

	if !r.KeyExists(key) {
		r.SetKey(key, "1", 1)
		return nil
	}

	err := r.RedisClient.IncrBy(r.Ctx, key, 1).Err()
	if err != nil {
		return err
	}

	if r.ShouldBlock(key) {
		k := fmt.Sprintf("blocked:%s", key)
		r.SetKey(k, "true", r.BlockedTimeInSeconds)
	}

	return nil
}

func (r *RedisConfig) IsBlocked(key string) bool {
	k := fmt.Sprintf("blocked:%s", key)

	if !r.KeyExists(k) {
		return false
	}

	val, _ := r.GetKey(k)

	return val != nil
}

func (r *RedisConfig) ShouldBlock(key string) bool {
	val, _ := r.GetKey(key)
	i, err := strconv.Atoi(*val)
	if err != nil {
		panic(err)
	}

	return i >= r.RequestsPerSecond
}
