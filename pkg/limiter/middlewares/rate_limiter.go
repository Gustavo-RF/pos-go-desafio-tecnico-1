package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/Gustavo-RF/desafio-tecnico-1/pkg/limiter"
)

type Response struct {
	Message string `json:"message"`
}

func RateLimiter(l *limiter.RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			success := l.Limiter(w, r)

			if success {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusBadGateway)
				response := Response{
					Message: "request failed. Check your configurations",
				}
				json.NewEncoder(w).Encode(response)
			}
		})
	}
}
