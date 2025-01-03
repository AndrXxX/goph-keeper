package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type ItemsController[T idItem] struct {
	Type string
	IF   sliceFetcher[T]
	IS   itemsStorage
	IC   itemConvertor[T]
}

func (c *ItemsController[T]) Update(w http.ResponseWriter, r *http.Request) {
	list, err := c.IF.FetchSlice(r.Body)
	if err != nil {
		logger.Log.Info("Update", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(enums.UserID).(uint)
	for _, item := range list {
		toSave, err := c.IC.ToModel(&item, userID)
		if err != nil {
			logger.Log.Info("c.IC.ToModel on ", zap.Error(err), zap.Any("item", item))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var cErr error
		if item.GetID() > 0 {
			exist, _ := c.IS.QueryOneById(r.Context(), item.GetID())
			if exist != nil && exist.UserID != userID {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			_, cErr = c.IS.Update(r.Context(), toSave)
		} else {
			_, cErr = c.IS.Insert(r.Context(), toSave)
		}

		if cErr != nil {
			logger.Log.Info("itemsService.Create", zap.Error(cErr))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (c *ItemsController[T]) Updates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value(enums.UserID).(uint)
	mList, err := c.IS.Query(r.Context(), &models.StoredItem{Type: datatypes.Text, UserID: userID})
	if err != nil {
		logger.Log.Info("c.IS.Query", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list := make([]T, len(mList))
	for i, item := range mList {
		rawItem, err := c.IC.ToEntity(&item)
		if err != nil {
			logger.Log.Info("c.IC.ToEntity", zap.Error(err), zap.Any("item", item))
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
		logger.Log.Info("w.Write on Updates TextItems", zap.Error(wErr))
	}
}
