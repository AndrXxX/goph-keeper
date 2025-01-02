package postgressql

import (
	"github.com/galeone/igor"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type Storage[T igor.DBModel] struct {
	db *igor.Database
}

func (s *Storage[T]) Find(m T) (T, error) {
	var resModel T
	err := s.db.Model(m).Where(m).Scan(resModel)
	if err != nil {
		logger.Log.Info("find model", zap.Error(err), zap.Any("model", m))
		return resModel, err
	}
	return resModel, nil
}

func (s *Storage[T]) Create(m T) (T, error) {
	err := s.db.Model(m).Create(m)
	if err != nil {
		logger.Log.Info("create model", zap.Error(err), zap.Any("model", m))
		return m, err
	}
	return m, nil
}

func (s *Storage[T]) Update(m T) error {
	err := s.db.Model(m).Updates(m)
	if err.Error != nil {
		logger.Log.Info("update model", zap.Error(err), zap.Any("model", m))
		return err
	}
	return nil
}

func (s *Storage[T]) List(m T) ([]T, error) {
	var l []T
	err := s.db.Model(m).Scan(&l)
	if err.Error != nil {
		logger.Log.Info("fetch list models", zap.Error(err), zap.Any("model", m))
		return l, err
	}
	return l, nil
}
