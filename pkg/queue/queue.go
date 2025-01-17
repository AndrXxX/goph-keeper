package queue

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AndrXxX/goph-keeper/pkg/logger"
)

type runner struct {
	si       time.Duration
	workers  []*worker
	jobs     chan Job
	running  atomic.Bool
	stopping atomic.Bool
	wg       sync.WaitGroup
}

func NewRunner(sleepInterval time.Duration) *runner {
	return &runner{
		si:      sleepInterval,
		workers: make([]*worker, 0),
	}
}

func (r *runner) SetWorkersCount(count int) *runner {
	r.workers = make([]*worker, count)
	return r
}

func (r *runner) AddJob(j Job) error {
	if !r.running.Load() {
		return errors.New("trying to add a queue before starting runner")
	}
	r.jobs <- j
	return nil
}

func (r *runner) Run(ctx context.Context) error {
	if r.running.Load() {
		return errors.New("already running")
	}
	r.jobs = make(chan Job)
	logger.Log.Info("Queue running")
	r.running.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	for {
		r.process(ctx)
		if r.stopping.Load() {
			cancel()
			r.stopping.Store(false)
			r.running.Store(false)
			return nil
		} else if len(r.workers) > 0 {
			r.wg.Wait()
		}
		time.Sleep(r.si)
	}
}

func (r *runner) process(ctx context.Context) {
	for i := range r.workers {
		r.wg.Add(1)
		go func(id int) {
			r.workers[i].Process(ctx, r.jobs)
			r.wg.Done()
		}(i)
	}
}

func (r *runner) Stop(ctx context.Context) error {
	select {
	default:
		logger.Log.Info("Queue shutting down")
		r.stopping.Store(true)
		for {
			if !r.running.Load() {
				break
			}
			time.Sleep(r.si)
		}
		logger.Log.Info("Queue stopped")
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
