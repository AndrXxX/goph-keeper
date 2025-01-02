package postgressql

import (
	"context"
	"fmt"

	"github.com/vingarcia/ksql"

	"github.com/AndrXxX/goph-keeper/pkg/storages/postgressql/models"
)

var storedItemsTable = ksql.NewTable("stored_items", "id")

type storedItemsStorage struct {
	db ksql.Provider
}

// Insert вставляет запись
func (s *storedItemsStorage) Insert(ctx context.Context, m *models.StoredItem) (*models.StoredItem, error) {
	err := s.db.Insert(ctx, storedItemsTable, m)
	if err != nil {
		return nil, fmt.Errorf("insert StoredItem %w", err)
	}
	return m, nil
}

// QueryOneById извлекает одну запись
func (s *storedItemsStorage) QueryOneById(ctx context.Context, m *models.StoredItem) (*models.StoredItem, error) {
	res := models.StoredItem{}
	err := s.db.QueryOne(ctx, &res, "FROM stored_items WHERE id = $1", m.ID)
	if err != nil {
		return nil, fmt.Errorf("queryOne StoredItem %w", err)
	}
	return &res, err
}

// Query извлекает несколько записей
func (s *storedItemsStorage) Query(ctx context.Context, m *models.StoredItem) ([]models.StoredItem, error) {
	var res []models.StoredItem
	q := "FROM stored_items WHERE type = $1 AND user_id = $2"
	params := []any{m.Type, m.UserID}
	if m.UpdatedAt != nil {
		q += " AND updated_at >= $3"
		params = append(params, *m.UpdatedAt)
	}
	err := s.db.Query(ctx, &res, q, params...)
	if err != nil {
		return nil, fmt.Errorf("query StoredItem %w", err)
	}
	return res, err
}

// Update обновляет запись
func (s *storedItemsStorage) Update(ctx context.Context, m *models.StoredItem) (*models.StoredItem, error) {
	err := s.db.Patch(ctx, storedItemsTable, m)
	if err != nil {
		return nil, fmt.Errorf("update StoredItem %w", err)
	}
	return m, nil
}
