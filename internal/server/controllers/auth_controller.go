package controllers

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type AuthController struct {
	US      usersStorage
	HG      hashGenerator
	TS      tokenService
	UF      fetcher[entities.User]
	KeyPath string
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	u, err := c.UF.Fetch(r.Body)
	if err != nil {
		logger.Log.Info("failed to fetchUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exist, err := c.US.QueryOne(r.Context(), u.Login)
	if err != nil {
		logger.Log.Info("Register QueryOne", zap.Error(err))
	}
	if exist != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	created, err := c.US.Insert(r.Context(), &models.User{Login: u.Login, Password: c.HG.Generate([]byte(u.Password))})
	if err != nil {
		logger.Log.Info("create user on register request", zap.Error(err))
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err := c.setAuthToken(w, created); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.writePublicKey(w)
	if err != nil {
		logger.Log.Info("failed to writePublicKey", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	u, err := c.UF.Fetch(r.Body)
	if err != nil {
		logger.Log.Info("fetchUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exist, err := c.US.QueryOne(r.Context(), u.Login)
	if err != nil {
		logger.Log.Info("Login QueryOne", zap.Error(err))
	}
	if exist == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if exist.Password != c.HG.Generate([]byte(u.Password)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := c.setAuthToken(w, exist); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.writePublicKey(w)
	if err != nil {
		logger.Log.Info("failed to writePublicKey", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) setAuthToken(w http.ResponseWriter, user *models.User) error {
	token, err := c.TS.Encrypt(user.ID)
	if err != nil {
		logger.Log.Info("failed to encrypt token on request", zap.Error(err))
		return err
	}
	w.Header().Set("Authorization", "Bearer "+token)
	http.SetCookie(w, &http.Cookie{Name: enums.AuthToken, Value: token})
	return err
}

func (c *AuthController) writePublicKey(w http.ResponseWriter) error {
	if c.KeyPath == "" {
		return nil
	}
	data, err := os.ReadFile(c.KeyPath)
	if err != nil {
		logger.Log.Info("failed to read cryptoKey", zap.Error(err), zap.String("keyPath", c.KeyPath))
		return nil
	}
	_, err = w.Write(data)
	return err
}
