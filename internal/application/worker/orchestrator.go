package worker

import (
	"fmt"

	"discordgo-template/pkg/logger"
)

type registeredWorker struct {
	name   string
	worker Worker
}

type Orchestrator struct {
	workers []registeredWorker
	logger  logger.Logger
}

func NewOrchestrator(logger logger.Logger) *Orchestrator {
	return &Orchestrator{
		workers: make([]registeredWorker, 0),
		logger:  logger,
	}
}

func (s *Orchestrator) RegisterWorker(name string, worker Worker) {
	for _, w := range s.workers {
		if w.name == name {
			s.logger.Warn("Worker %s already registered", name)
			return
		}
	}

	s.workers = append(s.workers, registeredWorker{name: name, worker: worker})
	s.logger.Info("Registered worker %s", name)
}

func (s *Orchestrator) StartAll() error {
	for _, w := range s.workers {
		if err := w.worker.Start(); err != nil {
			return fmt.Errorf("failed to start %s worker: %v", w.name, err)
		}
		s.logger.Info("Started worker %s", w.name)
	}
	return nil
}

func (s *Orchestrator) StopAll() error {
	for _, w := range s.workers {
		if err := w.worker.Stop(); err != nil {
			return fmt.Errorf("failed to stop %s worker: %v", w.name, err)
		}
		s.logger.Info("Stopped worker %s", w.name)
	}
	return nil
}
