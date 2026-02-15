package webhook

import (
	"sync"
	"time"
)

// LogStore stores delivery logs in memory
type LogStore struct {
	logs []DeliveryLog
	mu   sync.RWMutex
}

// NewLogStore creates a new log store
func NewLogStore() *LogStore {
	return &LogStore{
		logs: make([]DeliveryLog, 0, 1000),
	}
}

// Add adds a delivery log
func (s *LogStore) Add(log *DeliveryLog) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = append(s.logs, *log)

	// Keep only last 1000 logs
	if len(s.logs) > 1000 {
		s.logs = s.logs[len(s.logs)-1000:]
	}
}

// GetRecent returns recent logs
func (s *LogStore) GetRecent(limit int) []DeliveryLog {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > len(s.logs) {
		limit = len(s.logs)
	}

	start := len(s.logs) - limit
	if start < 0 {
		start = 0
	}

	result := make([]DeliveryLog, limit)
	copy(result, s.logs[start:])
	return result
}

// GetByChannel returns logs for a specific channel
func (s *LogStore) GetByChannel(channelID string, limit int) []DeliveryLog {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []DeliveryLog
	for i := len(s.logs) - 1; i >= 0 && len(result) < limit; i-- {
		if s.logs[i].ChannelID == channelID {
			result = append(result, s.logs[i])
		}
	}
	return result
}

// GetByStatus returns logs with a specific status
func (s *LogStore) GetByStatus(status DeliveryStatus, limit int) []DeliveryLog {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []DeliveryLog
	for i := len(s.logs) - 1; i >= 0 && len(result) < limit; i-- {
		if s.logs[i].Status == status {
			result = append(result, s.logs[i])
		}
	}
	return result
}

// Cleanup removes logs older than the specified duration
func (s *LogStore) Cleanup(maxAge time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	var filtered []DeliveryLog
	for _, log := range s.logs {
		if log.CreatedAt.After(cutoff) {
			filtered = append(filtered, log)
		}
	}
	s.logs = filtered
}
