package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := ""

		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}

		if tokenStr == "" {
			apiKey := r.URL.Query().Get("api_key")
			if apiKey != "" {
				k, err := h.svc.GetAPIKey(apiKey)
				if err == nil && (k.ExpiresAt == nil || k.ExpiresAt.After(time.Now())) {
					user, err := h.svc.GetUserByID(k.UserID)
					if err == nil {
						ctx := context.WithValue(r.Context(), "user", user)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
			writeError(w, 401, "unauthorized")
			return
		}

		secret := h.cfg.JWT.Secret
		if secret == "" {
			secret = os.Getenv("FOLIO_JWT_SECRET")
		}
		if secret == "" {
			secret = "change-me-in-production"
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			writeError(w, 401, "invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, 401, "invalid token")
			return
		}

		userID := uint(claims["user_id"].(float64))
		user, err := h.svc.GetUserByID(userID)
		if err != nil {
			writeError(w, 401, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
