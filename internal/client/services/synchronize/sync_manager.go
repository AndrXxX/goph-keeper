package synchronize

import (
	"errors"
	"fmt"

	"github.com/AndrXxX/goph-keeper/internal/client/services/synchronize/e"
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
	err := s.Sync(updates)
	if errors.Is(err, e.UnauthorizedError) {
		m.TR()
		err = s.Sync(updates)
	}
	return err
}
