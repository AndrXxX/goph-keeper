package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type AuthController struct {
	US userService
	HG hashGenerator
	TS tokenService
	UF userJSONRequestFetcher
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	u, err := c.UF.Fetch(r.Body)
	if err != nil {
		logger.Log.Info("failed to fetchUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exist, err := c.US.Find(&models.User{Login: u.Login})
	if exist != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	created, err := c.US.Create(&models.User{Login: u.Login, Password: c.HG.Generate([]byte(u.Password))})
	if err != nil {
		logger.Log.Info("failed to create user on register request", zap.Error(err))
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err := c.setAuthToken(w, created); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	u, err := c.UF.Fetch(r.Body)
	if err != nil {
		logger.Log.Info("failed to fetchUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exist, err := c.US.Find(&models.User{Login: u.Login})
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
	w.WriteHeader(http.StatusOK)
}

func (c *AuthController) setAuthToken(w http.ResponseWriter, user *models.User) error {
	token, err := c.TS.Encrypt(user.ID)
	if err != nil {
		logger.Log.Info("failed to encrypt token on request", zap.Error(err))
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: enums.AuthToken, Value: token})
	return err
}
