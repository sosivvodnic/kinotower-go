package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sosivvodnic/kinotower-go/internal/core/auth"
	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func UserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(userIDKey)
	id, ok := v.(int)
	return id, ok
}

func RequireAuth(mgr *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" {
				httpjson.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httpjson.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims, err := mgr.Parse(parts[1])
			if err != nil {
				httpjson.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

