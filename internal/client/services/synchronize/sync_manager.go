package synchronize

import "fmt"

type SyncManager struct {
	Synchronizers map[string]Synchronizer
}

func (m SyncManager) Sync(dataType string, updates []any) error {
	if s, ok := m.Synchronizers[dataType]; ok {
		return s.Sync(updates)
	}
	return fmt.Errorf("unknown data type: %s", dataType)
}
