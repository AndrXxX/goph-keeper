package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/internal/enums/contenttypes"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type ItemsController[T idItem] struct {
	Type      string
	Fetcher   sliceFetcher[T]
	Storage   itemsStorage
	Convertor itemConvertor[T]
}

func (c *ItemsController[T]) StoreUpdates(w http.ResponseWriter, r *http.Request) {
	list, err := c.Fetcher.FetchSlice(r.Body)
	if err != nil {
		logger.Log.Info("StoreUpdates", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(enums.UserID).(uint)
	for _, item := range list {
		toSave, err := c.Convertor.ToModel(&item, userID)
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
	}
	w.WriteHeader(http.StatusOK)
}

func (c *ItemsController[T]) FetchUpdates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contenttypes.ApplicationJSON)
	userID := r.Context().Value(enums.UserID).(uint)
	mList, err := c.Storage.Query(r.Context(), &models.StoredItem{Type: c.Type, UserID: userID})
	if err != nil {
		logger.Log.Info("c.Storage.Query", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list := make([]T, len(mList))
	for i, item := range mList {
		rawItem, err := c.Convertor.ToEntity(&item)
		if err != nil {
			logger.Log.Info("c.Convertor.ToEntity", zap.Error(err), zap.Any("item", item))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		list[i] = *rawItem
	}
	res, mErr := json.Marshal(list)
	if mErr != nil {
		logger.Log.Info("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, wErr := w.Write(res)
	if wErr != nil {
		logger.Log.Info("w.Write on FetchUpdates TextItems", zap.Error(wErr))
	}
}
