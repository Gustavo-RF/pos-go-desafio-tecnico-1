package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Gustavo-RF/desafio-tecnico-1/configs"
	"github.com/Gustavo-RF/desafio-tecnico-1/internal/handler"
	"github.com/Gustavo-RF/desafio-tecnico-1/internal/infra/redis"
	"github.com/Gustavo-RF/desafio-tecnico-1/pkg/limiter"
	"github.com/Gustavo-RF/desafio-tecnico-1/pkg/limiter/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	// carrega as configs
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// inicia o redis
	ctx := context.Background()
	redisClient := redis.NewRedisClient(
		ctx,
		configs.RedisAddress,
		configs.RedisPort,
		configs.RedisPassword,
		configs.RateLimiter,
		configs.RequestsPerSecond,
		configs.BlockedTimeInSeconds,
	)

	if redisClient == nil {
		panic(errors.New("error while load redis"))
	}

	// inicia o rate limiter
	rateLimiter := limiter.NewRateLimiterConfig(redisClient, configs.RateLimiter, configs.RequestsPerSecond, configs.BlockedTimeInSeconds)

	// inicia o server
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middlewares.RateLimiter(rateLimiter))

	r.Get("/", handler.Handle)

	fmt.Printf("Server started at %s...\n", configs.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.Port), r)
}
