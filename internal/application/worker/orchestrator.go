package worker

import (
	"fmt"

	"discordgo-template/pkg/logger"
)

// question: is it a smell to need to use an ordered map?
// todo: make workers an ordered map
type Orchestrator struct {
	workers map[string]Worker
	logger  logger.Logger
}

func NewOrchestrator(logger logger.Logger) *Orchestrator {
	return &Orchestrator{
		workers: make(map[string]Worker),
		logger:  logger,
	}
}

func (s *Orchestrator) RegisterWorker(name string, worker Worker) {
	if _, ok := s.workers[name]; ok {
		s.logger.Warn("Worker %s already registered", name)
		return
	}

	s.workers[name] = worker
	s.logger.Info("Registered worker %s", name)
}

func (s *Orchestrator) StartAll() error {
	for name, worker := range s.workers {
		if err := worker.Start(); err != nil {
			return fmt.Errorf("failed to start %s worker: %v", name, err)
		}
		s.logger.Info("Started worker %s", name)
	}
	return nil
}

func (s *Orchestrator) StopAll() error {
	for name, worker := range s.workers {
		if err := worker.Stop(); err != nil {
			return fmt.Errorf("failed to stop %s worker: %v", name, err)
		}
		s.logger.Info("Stopped worker %s", name)
	}
	return nil
}
