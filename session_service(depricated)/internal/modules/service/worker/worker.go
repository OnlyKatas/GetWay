package worker

import (
	"context"
	"log"
	"session/internal/models"
)

type Pool struct {
	tasks    chan models.Session
	workers  []*Worker
	stopChan chan struct{}
}

type Worker struct {
	id     int
	taskCh chan models.Session
	stopCh chan struct{}
	ctx    context.Context
}

func NewPool(size int) *Pool {
	p := &Pool{
		tasks:    make(chan models.Session),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < size; i++ {
		worker := &Worker{
			id:     i,
			taskCh: p.tasks,
			stopCh: p.stopChan,
		}
		p.workers = append(p.workers, worker)
	}

	return p
}

func (p *Pool) Start(ctx context.Context) {
	for _, worker := range p.workers {
		go worker.Start(ctx)
	}
}

func (p *Pool) AddTask(task models.Session) {
	p.tasks <- task
}

func (p *Pool) Stop() {
	close(p.tasks)
	close(p.stopChan)
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case task, ok := <-w.taskCh:
			if !ok {
				return
			}
			w.Work(ctx, &task)
		case <-w.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) Work(ctx context.Context, s *models.Session) {
	log.Printf("Worker %d processing session: %v", w.id, s)
}
