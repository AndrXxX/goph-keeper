package controllers

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/internal/enums/contenttypes"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type FilesController struct {
	Storage itemsStorage
	FS      fileStorage
}

func (c *FilesController) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contenttypes.FormUrlEncoded)

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		logger.Log.Info("ParseMultipartForm", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := r.Form.Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	objId, err := uuid.FromBytes([]byte(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, _ := c.Storage.QueryOneById(r.Context(), objId)
	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userID := r.Context().Value(enums.UserID).(uint)
	if m.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if c.FS.IsExist(objId) {
		w.WriteHeader(http.StatusOK)
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		logger.Log.Info("FormFile", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)

	err = c.FS.Store(f, objId)
	if err != nil {
		logger.Log.Info("Store", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FilesController) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	objId, err := uuid.FromBytes([]byte(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, _ := c.Storage.QueryOneById(r.Context(), objId)
	if m == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userID := r.Context().Value(enums.UserID).(uint)
	if m.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !c.FS.IsExist(objId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := c.FS.Get(objId)
	if err != nil {
		logger.Log.Info("FS.Get", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, file)
	if err != nil {
		logger.Log.Error("io.Copy", zap.Error(err))
		return
	}
}
