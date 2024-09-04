package middleware

import (
	"JwtWithGo/pkg/controllers"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

func IsAuthorized(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}

			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := authHeaderParts[1]
			claims := &controllers.Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return controllers.JwtKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if claims.ExpiresAt < time.Now().Unix() {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
