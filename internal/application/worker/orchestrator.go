package worker

import (
	"fmt"

	"discordgo-skeleton/pkg/logger"
)

type ManagedWorker struct {
	id     string
	worker Worker
}

type Orchestrator struct {
	workers    []*ManagedWorker
	registered map[string]struct{}
	logger     logger.Logger
}

func NewOrchestrator(logger logger.Logger) *Orchestrator {
	return &Orchestrator{
		workers:    make([]*ManagedWorker, 0),
		registered: make(map[string]struct{}),
		logger:     logger,
	}
}

func (s *Orchestrator) RegisterWorker(id string, worker Worker) error {
	if _, ok := s.registered[id]; ok {
		return fmt.Errorf("Worker %s already registered", id)
	}
	s.registered[id] = struct{}{}

	mw := &ManagedWorker{
		id:     id,
		worker: worker,
	}
	s.workers = append(s.workers, mw)
	s.logger.Info("Registered worker %s", id)
	return nil
}

func (s *Orchestrator) StartAll() error {
	for _, mw := range s.workers {
		if err := mw.worker.Start(); err != nil {
			return err
		}
		s.logger.Info("Started worker %s", mw.id)
	}
	return nil
}

func (s *Orchestrator) StopAll() error {
	for _, mw := range s.workers {
		if err := mw.worker.Stop(); err != nil {
			return err
		}
		s.logger.Info("Stopped worker %s", mw.id)
	}
	return nil
}
