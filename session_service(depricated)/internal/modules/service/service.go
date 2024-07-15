package service

import (
	"context"
	"log"
	"session/internal/models"
	workers "session/internal/modules/service/worker"
	"time"
)

type Service struct {
	Payer      Payer
	workerPool *workers.Pool
}

type Payer interface {
	GetSessionStatus(ctx context.Context) ([]models.Session, error)
}

func NewService(client Payer, poolSize int) *Service {
	return &Service{
		Payer:      client,
		workerPool: workers.NewPool(poolSize),
	}
}

func (s *Service) WorkerPool(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sessions, err := s.Payer.GetSessionStatus(ctx)
			if err != nil {
				log.Println("Error getting session status:", err)
				continue
			}

			// Start the worker pool
			s.workerPool.Start(ctx)

			// Send sessions to the worker pool
			for _, session := range sessions {
				select {
				case <-ctx.Done():
					log.Println("Worker pool stopped.")
					return
				default:
					s.workerPool.AddTask(session)
				}
			}

		case <-ctx.Done():
			log.Println("Worker pool stopped.")
			return
		}
	}
}
