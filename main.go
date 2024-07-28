package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Gustavo-RF/desafio-tecnico-1/configs"
	"github.com/Gustavo-RF/desafio-tecnico-1/internal/handler"
	"github.com/Gustavo-RF/desafio-tecnico-1/internal/infra/redisconfig"
	"github.com/Gustavo-RF/desafio-tecnico-1/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// PEGA A REQUEST
	// VER SE POSSUI A KEY
	// SE POSSUIR, VALIDAR PELA KEY
	// BLOQUEAR CASO ULTRAPASSE O CONFIGURADO NO TEMPO

	// carrega as configs
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// inicia o redis
	ctx := context.Background()
	redisClient := redisconfig.NewRedisClient(
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

	// inicia o server
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middlewares.RateLimiter)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler.Handle(*redisClient, w, r)
	})

	fmt.Printf("Server started at %s...\n", configs.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.Port), r)
}
