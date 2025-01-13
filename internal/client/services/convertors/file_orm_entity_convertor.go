package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type FileOrmEntityConvertor struct {
}

func (c FileOrmEntityConvertor) Convert(m *ormmodels.FileItem) *entities.FileItem {
	if m == nil {
		return nil
	}
	return &entities.FileItem{
		StoredItem: *ItemOrmEntityConvertor{}.Convert(&m.StoredItem),
		Name:       m.Name,
	}
}
