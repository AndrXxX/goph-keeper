package app

import (
	"context"

	"github.com/AndrXxX/goph-keeper/pkg/queue"
)

type queueRunner interface {
	Run(context.Context) error
	Stop(context.Context) error
	AddJob(queue.Job) error
}

type syncManager interface {
	Sync(dataType string, updates []any) error
}
