package middlewares

import (
	"fmt"
	"net/http"
)

func RateLimiter(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Rate limiter")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
