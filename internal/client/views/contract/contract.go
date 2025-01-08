package contract

type SyncManager interface {
	Sync(dataType string, updates []any) error
}
