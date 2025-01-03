package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/internal/enums"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/entities/values"
	"github.com/AndrXxX/goph-keeper/pkg/logger"
	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

type LoginPassItemsController struct {
	IF sliceFetcher[entities.LoginPassItem]
	IS itemsStorage
}

func (c *LoginPassItemsController) Update(w http.ResponseWriter, r *http.Request) {
	list, err := c.IF.FetchSlice(r.Body)
	if err != nil {
		logger.Log.Info("failed to Update", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(enums.UserID).(uint)
	for _, item := range list {
		val, err := json.Marshal(values.LoginPassValue{Login: item.Login, Password: item.Password})
		if err != nil {
			logger.Log.Info("json.Marshal", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		toSave := models.StoredItem{
			ID:          item.ID,
			Type:        datatypes.LoginPass,
			Description: item.Description,
			Value:       string(val),
			UserID:      userID,
		}
		var cErr error
		if item.ID > 0 {
			exist, _ := c.IS.QueryOneById(r.Context(), &models.StoredItem{ID: item.ID})
			if exist != nil && exist.UserID != userID {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			_, cErr = c.IS.Update(r.Context(), &toSave)
		} else {
			_, cErr = c.IS.Insert(r.Context(), &toSave)
		}

		if cErr != nil {
			logger.Log.Info("itemsService.Create", zap.Error(cErr))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (c *LoginPassItemsController) Updates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value(enums.UserID).(uint)
	mList, err := c.IS.Query(r.Context(), &models.StoredItem{Type: datatypes.LoginPass, UserID: userID})
	if err != nil {
		logger.Log.Info("c.IS.Query", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	list := make([]entities.LoginPassItem, len(mList))
	for i, item := range mList {
		var rawItem entities.LoginPassItem
		err := json.Unmarshal([]byte(item.Value), &rawItem)
		if err != nil {
			logger.Log.Info("json.Unmarshal", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rawItem.StoredItem = entities.StoredItem{ID: item.ID, Description: item.Description}
		list[i] = rawItem
	}
	res, mErr := json.Marshal(list)
	if mErr != nil {
		logger.Log.Info("json.Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(res)
	if wErr != nil {
		logger.Log.Info("w.Write", zap.Error(wErr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
