package jobs

type UploadItemsUpdatesJob struct {
	Type        string
	Items       []any
	SyncManager syncManager
}

func (j *UploadItemsUpdatesJob) Execute() error {
	return j.SyncManager.Sync(j.Type, j.Items)
}
