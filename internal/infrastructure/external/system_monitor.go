package external

import (
	"time"

	"discordgo-template/internal/application/ports"

	"discordgo-template/pkg/httpx"
)

var _ ports.SystemMonitor = (*SystemMonitor)(nil)

type SystemMonitor struct {
	client    httpx.Client
	startTime int64
}

func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{
		client: httpx.NewRetryClient(
			httpx.WithRetryPolicy(httpx.NoRetry()),
			httpx.WithBackoffPolicy(httpx.NoBackoff()),
		),
		startTime: time.Now().Unix(),
	}
}

func (p *SystemMonitor) StartTime() int64 {
	return p.startTime
}
