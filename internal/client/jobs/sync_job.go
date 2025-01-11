package jobs

type SyncJob struct {
	Type        string
	SyncManager syncManager
}

func (j *SyncJob) Execute() error {
	return j.SyncManager.Sync(j.Type, []any{})
}
