package middlewares

import (
	"fmt"
	"net/http"
)

func CustomMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Testes middleware")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
