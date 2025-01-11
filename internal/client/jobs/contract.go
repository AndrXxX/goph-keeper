package jobs

import (
	"github.com/AndrXxX/goph-keeper/pkg/queue"
)

type syncManager interface {
	Sync(dataType string, updates []any) error
}

type queueRunner interface {
	AddJob(queue.Job) error
}
