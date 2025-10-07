package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"workoutpal/src/util/constants"

	"github.com/go-chi/chi/v5"
)

func IdMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			stringID := chi.URLParam(r, constants.ID_KEY)
			if stringID == "" {
				http.Error(w, "id not provided in path", http.StatusUnauthorized)
				return
			}

			id, err := strconv.ParseInt(stringID, 10, 64)
			if err != nil {
				fmt.Println("provided ", id, stringID)
				http.Error(w, "id provided is not a valid number", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), constants.ID_KEY, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
