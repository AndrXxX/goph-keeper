package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type FileEntityOrmConvertor struct {
}

func (c FileEntityOrmConvertor) Convert(e *entities.FileItem) *ormmodels.FileItem {
	if e == nil {
		return nil
	}
	return &ormmodels.FileItem{
		StoredItem: *ItemEntityOrmConvertor{}.Convert(&e.StoredItem),
		Name:       e.Name,
	}
}
