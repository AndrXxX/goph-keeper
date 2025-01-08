package messages

import "github.com/AndrXxX/goph-keeper/internal/client/entities"

type AddNote struct {
	Item *entities.NoteItem
}
