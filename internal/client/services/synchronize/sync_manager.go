package synchronize

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/syncerr"
)

type SyncManager struct {
	Synchronizers map[string]Synchronizer
	TR            tokenRefresher
}

func (m SyncManager) Sync(dataType string, updates []any) error {
	s, ok := m.Synchronizers[dataType]
	if !ok {
		return fmt.Errorf("unknown data type: %s", dataType)
	}
	for _, item := range updates {
		if _, err := govalidator.ValidateStruct(item); err != nil {
			return err
		}
	}
	err := s.Sync(updates)
	if errors.Is(err, syncerr.UnauthorizedError) {
		m.TR()
		err = s.Sync(updates)
	}
	return err
}
