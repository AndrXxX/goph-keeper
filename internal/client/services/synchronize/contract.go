package synchronize

type Synchronizer interface {
	Sync(updates []any) error
}
