package handler

import (
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
	ip, err := getIP(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiKey := r.Header.Get("Api-Key")

	if apiKey != "" {
		fmt.Printf("Has key: %s\n", apiKey)
	}

	fmt.Printf("Has ip: %s\n", ip)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ip))
}
