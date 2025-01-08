package synchronizers

import (
	"fmt"
	"net/http"

	"github.com/AndrXxX/goph-keeper/internal/client/entities"
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/convertors"
	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
)

type NotesSynchronizer struct {
	LC convertors.ListConvertor[entities.NoteItem]
	L  loader[entities.NoteItem]
	s  interfaces.Storage[entities.NoteItem]
}

func (s *NotesSynchronizer) Sync(updates []any) error {
	list, cErr := s.LC.Convert(updates)
	if cErr != nil {
		return fmt.Errorf("convert note updates: %w", cErr)
	}
	uErr := s.L.Upload(datatypes.Notes, list)
	if uErr != nil {
		return fmt.Errorf("upload note updates: %w", cErr)
	}
	dErr := s.Download()
	if dErr != nil {
		return fmt.Errorf("download note updates: %w", dErr)
	}
	return nil
}

func (s *NotesSynchronizer) Download() error {
	code, list := s.L.Download(datatypes.Notes)
	if code != http.StatusOK {
		return fmt.Errorf("download notes - unexpected code: %v", code)
	}
	for i := range list {
		exist := s.s.Find(&list[i])
		if exist == nil {
			_, err := s.s.Create(exist)
			if err != nil {
				return fmt.Errorf("create note item: %w", err)
			}
			continue
		}
		exist.Text = list[i].Text
		exist.StoredItem = list[i].StoredItem
		err := s.s.Update(exist)
		if err != nil {
			return fmt.Errorf("update note item: %w", err)
		}
	}
	return nil
}
