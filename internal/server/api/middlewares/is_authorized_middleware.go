package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type isAuthorized struct {
	ts tokenService
}

func (m *isAuthorized) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(enums.AuthToken)
		if err != nil {
			logger.Log.Info("failed to get auth token cookie", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := m.ts.Decrypt(token.Value)
		if err != nil {
			logger.Log.Info("failed to decrypt token from cookie", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, enums.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsAuthorized(ts tokenService) *isAuthorized {
	return &isAuthorized{ts}
}
