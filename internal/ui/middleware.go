package ui

import (
	"context"
	"net/http"
)

type ctxKey string

const layoutKey ctxKey = "layout"

func Layout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		layout := "base"

		if r.URL.Path == "/login" {
			layout = "auth"
		}

		ctx := context.WithValue(r.Context(), layoutKey, layout)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
