package usecase

import (
	"time"

	"discordgo-template/internal/application/ports"
)

type SystemService struct {
	monitor ports.SystemMonitor
}

func NewSystemService(monitor ports.SystemMonitor) *SystemService {
	return &SystemService{
		monitor: monitor,
	}
}

// --- Get Status ---
type GetStatusRequest struct{}

type GetStatusResult struct {
	StartTime     int64
	UptimeDays    int64
	UptimeHours   int64
	UptimeMinutes int64
	UptimeSeconds int64
}

func (s *SystemService) GetStatus(req GetStatusRequest) (*GetStatusResult, error) {
	startTime := s.monitor.StartTime()
	diff := time.Now().Unix() - startTime

	res := &GetStatusResult{
		StartTime:     s.monitor.StartTime(),
		UptimeDays:    diff / 86400,
		UptimeHours:   (diff / 3600) % 24,
		UptimeMinutes: (diff / 60) % 60,
		UptimeSeconds: diff % 60,
	}
	return res, nil
}
