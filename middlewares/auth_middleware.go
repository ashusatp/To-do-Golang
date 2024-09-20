package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/ashusatp/todo/models"
	"github.com/dgrijalva/jwt-go"
)

type key int

const UserContextKey key = iota

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.Split(authHeader, " ")[1]

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return models.JwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
