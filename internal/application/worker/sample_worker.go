package worker

import (
	"fmt"
	"sync"
	"time"

	"discordgo-template/pkg/logger"
)

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type realTicker struct {
	*time.Ticker
}

func (t realTicker) C() <-chan time.Time {
	return t.Ticker.C
}

type TickerFactory func(time.Duration) Ticker

type SampleWorker struct {
	interval      time.Duration
	logger        logger.Logger
	tickerFactory TickerFactory

	mu       sync.Mutex
	running  bool
	stopping bool
	stop     chan struct{}
	done     chan struct{}
}

func NewSampleWorker(interval time.Duration, logger logger.Logger, options ...func(*SampleWorker)) *SampleWorker {
	worker := &SampleWorker{
		interval: interval,
		logger:   logger,
		tickerFactory: func(interval time.Duration) Ticker {
			return realTicker{Ticker: time.NewTicker(interval)}
		},
	}
	for _, option := range options {
		option(worker)
	}
	return worker
}

func WithTickerFactory(factory TickerFactory) func(*SampleWorker) {
	return func(worker *SampleWorker) {
		worker.tickerFactory = factory
	}
}

func (w *SampleWorker) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.running {
		return nil
	}
	if w.stopping {
		return fmt.Errorf("sample worker is stopping")
	}
	if w.interval <= 0 {
		return fmt.Errorf("sample worker interval must be greater than zero")
	}

	ticker := w.tickerFactory(w.interval)
	w.stop = make(chan struct{})
	w.done = make(chan struct{})
	w.running = true
	go w.run(ticker, w.stop, w.done)
	return nil
}

func (w *SampleWorker) run(ticker Ticker, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	defer ticker.Stop()
	for {
		select {
		case tick := <-ticker.C():
			w.logger.Info("sample worker tick: %v", tick)
		case <-stop:
			return
		}
	}
}

func (w *SampleWorker) Stop() error {
	w.mu.Lock()
	if !w.running {
		done := w.done
		w.mu.Unlock()
		if done != nil {
			<-done
		}
		return nil
	}
	w.running = false
	w.stopping = true
	stop := w.stop
	done := w.done
	close(stop)
	w.mu.Unlock()

	<-done

	w.mu.Lock()
	w.stopping = false
	w.stop = nil
	w.done = nil
	w.mu.Unlock()
	return nil
}
