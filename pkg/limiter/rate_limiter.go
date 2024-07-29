package limiter

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/Gustavo-RF/desafio-tecnico-1/internal/entities"
	"github.com/Gustavo-RF/desafio-tecnico-1/pkg/limiter/interfaces"
)

type RateLimiterConfig struct {
	RateLimiter          interfaces.RateLimiterInterface
	RateLimiterType      string
	RequestsPerSecond    int
	BlockedTimeInSeconds int
}

func NewRateLimiterConfig(rate interfaces.RateLimiterInterface, rlType string, rps, btis int) *RateLimiterConfig {
	return &RateLimiterConfig{
		RateLimiter:          rate,
		RateLimiterType:      rlType,
		RequestsPerSecond:    rps,
		BlockedTimeInSeconds: btis,
	}
}

func (r *RateLimiterConfig) SetKey(key, value string, expireTime int) error {
	return r.RateLimiter.SetKey(key, value, expireTime)
}

func (r *RateLimiterConfig) GetKey(key string) (*string, error) {
	return r.RateLimiter.GetKey(key)
}

func (r *RateLimiterConfig) IncKey(key string) error {
	return r.RateLimiter.IncKey(key)
}

func (r *RateLimiterConfig) KeyExists(key string) bool {
	return r.RateLimiter.KeyExists(key)
}

func (r *RateLimiterConfig) IsBlocked(key string) bool {
	return r.RateLimiter.IsBlocked(key)
}

func (r *RateLimiterConfig) ShouldBlock(key string) bool {
	return r.RateLimiter.ShouldBlock(key)
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

func (r *RateLimiterConfig) Limiter(w http.ResponseWriter, req *http.Request) bool {
	var key string

	if r.RateLimiterType == "IP" {
		ip, err := getIP(req)
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			response := entities.Response{
				Message: "error while get ip",
			}
			json.NewEncoder(w).Encode(response)
			return false
		}

		key = ip
	} else if r.RateLimiterType == "TOKEN" {
		apiKey := req.Header.Get("Api-Key")

		if apiKey == "" {
			w.WriteHeader(http.StatusTooManyRequests)
			response := entities.Response{
				Message: "api key header is required when rate limiter is configured for token",
			}
			json.NewEncoder(w).Encode(response)
			return false
		}

		key = apiKey
	}

	if r.IsBlocked(key) {
		w.WriteHeader(http.StatusTooManyRequests)
		response := entities.Response{
			Message: "you have reached the maximum number of requests or actions allowed within a certain time frame",
		}
		json.NewEncoder(w).Encode(response)
		return false
	}

	r.IncKey(key)
	return true
}
