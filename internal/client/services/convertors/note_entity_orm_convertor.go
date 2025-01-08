package convertors

import (
	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/ormmodels"
)

type NoteEntityOrmConvertor struct {
}

func (c NoteEntityOrmConvertor) Convert(e *entities.NoteItem) *ormmodels.NoteItem {
	return &ormmodels.NoteItem{
		StoredItem: *ItemEntityOrmConvertor{}.Convert(&e.StoredItem),
		Text:       e.Text,
	}
}
