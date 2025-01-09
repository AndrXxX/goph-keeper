package queue

import (
	"context"

	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type worker struct {
}

func (w *worker) Process(ctx context.Context, jobs <-chan Job) {
	select {
	// if context was canceled.
	case <-ctx.Done():
		return
	// if job received.
	case job := <-jobs:
		err := job.Execute()
		if err != nil {
			logger.Log.Error("failed to execute runner job", zap.Error(err), zap.Any("job", job))
		}
	}
}
