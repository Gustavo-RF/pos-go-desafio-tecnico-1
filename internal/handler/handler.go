package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Gustavo-RF/desafio-tecnico-1/internal/infra/redisconfig"
)

type Handler struct {
	RedisConfig redisconfig.RedisConfig
}

type Response struct {
	Message string `json:"message"`
}

func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}

func Handle(rc redisconfig.RedisConfig, w http.ResponseWriter, r *http.Request) {

	var key string

	if rc.RateLimiter == "IP" {
		ip, err := getIP(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		key = ip
	} else if rc.RateLimiter == "TOKEN" {
		apiKey := r.Header.Get("Api-Key")

		if apiKey == "" {
			w.WriteHeader(http.StatusBadGateway)
			response := Response{
				Message: "api key header is required when rate limiter is configured for token",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		key = apiKey
	}

	fmt.Printf("Key: %s\n", key)
	if rc.Blocked(key) {
		w.WriteHeader(http.StatusTooManyRequests)
		response := Response{
			Message: "you have reached the maximum number of requests or actions allowed within a certain time frame",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	rc.IncKey(key)

	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "success",
	}
	json.NewEncoder(w).Encode(response)
}
