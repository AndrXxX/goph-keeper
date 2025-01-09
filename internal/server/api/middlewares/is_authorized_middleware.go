package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type isAuthorized struct {
	ts tokenService
}

func (m *isAuthorized) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.getTokenFromHeaders(r)
		if token == "" {
			t, err := m.getTokenFromCookie(r)
			if err != nil {
				logger.Log.Info("get auth token cookie", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token = t
		}

		userID, err := m.ts.Decrypt(token)
		if err != nil {
			logger.Log.Info("decrypt token", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, enums.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *isAuthorized) getTokenFromCookie(r *http.Request) (string, error) {
	t, err := r.Cookie(enums.AuthToken)
	if err != nil {
		return "", fmt.Errorf("getTokenFromCookie %w", err)
	}
	return t.Value, nil
}

func (m *isAuthorized) getTokenFromHeaders(r *http.Request) string {
	if raw := r.Header.Get("Authorization"); raw != "" {
		vals := strings.Split(raw, " ")
		if len(vals) == 2 {
			return vals[1]
		}
	}
	return ""
}

func IsAuthorized(ts tokenService) *isAuthorized {
	return &isAuthorized{ts}
}
