package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/config"
	"github.com/farinmark/Backend-trainee-assignment-spring-2025/internal/service"

	"github.com/dgrijalva/jwt-go"
)

type ctxKey string

const UserKey ctxKey = "user"

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func Authenticate(svc *service.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			tokenString := strings.TrimPrefix(auth, "Bearer ")
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.Load().JWTSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(*Claims)
			ctx := context.WithValue(r.Context(), UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
