package auth

import (
	"context"
	"net/http"
	"strings"
)

type Middleware struct {
	jwt *Jwt
}

func NewMiddleware(jwt *Jwt) *Middleware {
	return &Middleware{jwt: jwt}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.jwt.ParseToken(rawToken)
		if err != nil {
			// TODO: return a proper error response
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUserID, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
