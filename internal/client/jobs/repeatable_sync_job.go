package jobs

import (
	"time"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type RepeatableSyncJob struct {
	Type           string
	SyncManager    syncManager
	QR             queueRunner
	RepeatInterval time.Duration
}

func (j *RepeatableSyncJob) Execute() error {
	res := j.SyncManager.Sync(j.Type, []any{})
	go func() {
		time.Sleep(j.RepeatInterval)
		err := j.QR.AddJob(j)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}()
	return res
}
