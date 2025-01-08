package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type NoteOrmEntityConvertor struct {
}

func (c NoteOrmEntityConvertor) Convert(m *ormmodels.NoteItem) *entities.NoteItem {
	if m == nil {
		return nil
	}
	return &entities.NoteItem{
		StoredItem: *ItemOrmEntityConvertor{}.Convert(&m.StoredItem),
		Text:       m.Text,
	}
}
