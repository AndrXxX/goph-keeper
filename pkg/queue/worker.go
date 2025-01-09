package queue

import (
	"go.uber.org/zap"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type worker struct {
}

func (w *worker) Process(jobs <-chan Job) {
	for job := range jobs {
		err := job.Execute()
		if err != nil {
			logger.Log.Error("failed to execute runner job", zap.Error(err), zap.Any("job", job))
		}
	}
}
