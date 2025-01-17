package controllers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/internal/enums/contenttypes"
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type FilesController struct {
	Storage   itemsStorage
	FS        fileStorage
	FF        fetcher[entities.FileItem]
	Convertor itemConvertor[entities.FileItem]
}

func (c *FilesController) Update(w http.ResponseWriter, r *http.Request) {
	item, err := c.FF.Fetch(r.Body)
	if err != nil {
		logger.Log.Info("StoreUpdates", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(enums.UserID).(uint)
	toSave, err := c.Convertor.ToModel(item, userID)
	if err != nil {
		logger.Log.Info("c.Convertor.ToModel on ", zap.Error(err), zap.Any("item", item))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var cErr error

	if exist, _ := c.Storage.QueryOneById(r.Context(), item.GetID()); exist != nil {
		if exist.UserID != userID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, cErr = c.Storage.Update(r.Context(), toSave)
	} else {
		_, cErr = c.Storage.Insert(r.Context(), toSave)
	}

	if cErr != nil {
		logger.Log.Info("itemsService.Create", zap.Error(cErr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(toSave.ID.String()))
}

func (c *FilesController) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contenttypes.OctetStream)
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	objId, err := uuid.ParseBytes([]byte(id))
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

	err = c.FS.Store(r.Body, objId)
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
	objId, err := uuid.ParseBytes([]byte(id))
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
	w.Header().Set("Content-Type", contenttypes.OctetStream)
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		logger.Log.Error("io.Copy", zap.Error(err))
		return
	}
}
