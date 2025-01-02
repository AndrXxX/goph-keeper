package postgressql

import (
	"time"

	"github.com/galeone/igor"
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type Storage[T igor.DBModel] struct {
	db *igor.Database
}

func (s *Storage[T]) Find(m *T) (*T, error) {
	var res []T
	_ = s.db.Model(*new(T)).Where(*m).Scan(&res)
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (s *Storage[T]) Create(m *T) (*T, error) {
	err := s.db.Create(*m)
	if err != nil {
		logger.Log.Info("create model", zap.Error(err), zap.Any("model", *m))
		return m, err
	}
	return m, nil
}

func (s *Storage[T]) Update(m *T) error {
	err := s.db.Updates(*m)
	if err != nil {
		logger.Log.Info("update model", zap.Error(err), zap.Any("model", *m))
		return err
	}
	return nil
}

func (s *Storage[T]) List(m *T, updatedFrom *time.Time) ([]T, error) {
	var l []T
	q := s.db.Model(*m)
	if updatedFrom != nil {
		q = q.Where("updated_at > ?", updatedFrom)
	}
	err := q.Scan(&l)
	if err != nil {
		logger.Log.Info("fetch list models", zap.Error(err), zap.Any("model", m))
		return nil, err
	}
	return l, nil
}
