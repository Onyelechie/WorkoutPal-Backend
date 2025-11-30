package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"workoutpal/src/util/constants"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const claimsCtxKey ctxKey = iota

func ClaimsFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(claimsCtxKey).(jwt.MapClaims)
	return claims, ok
}

func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	if os.Getenv("APP_ENV") == "test" {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	if len(secret) == 0 {
		if env := os.Getenv("JWT_SECRET"); env != "" {
			secret = []byte(env)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			
			// First try to get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				// Fall back to cookie
				cookie, err := r.Cookie("access_token")
				if err != nil || cookie.Value == "" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(map[string]string{"error": "missing auth token"})
					return
				}
				tokenString = cookie.Value
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid or expired token"})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid token claims"})
				return
			}

			userIDFloat, _ := claims["sub"].(float64)
			userID := int64(userIDFloat)

			ctx := context.WithValue(r.Context(), claimsCtxKey, claims)
			ctx = context.WithValue(ctx, constants.USER_ID_KEY, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
