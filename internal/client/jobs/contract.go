package jobs

type syncManager interface {
	Sync(dataType string, updates []any) error
}
