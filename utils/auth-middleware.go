package utils

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "authenticated", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
