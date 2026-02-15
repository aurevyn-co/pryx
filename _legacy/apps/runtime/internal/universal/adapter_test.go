package universal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConnectionMetricsStructure tests ConnectionMetrics struct
func TestConnectionMetricsStructure(t *testing.T) {
	metrics := ConnectionMetrics{
		TotalConnections:  10,
		ActiveConnections: 5,
		MessagesSent:      100,
		MessagesReceived:  150,
		ErrorsTotal:       2,
		BytesSent:         1024,
		BytesReceived:     2048,
		LastActivity:      time.Now(),
		ProtocolStats: map[string]int64{
			"websocket": 5,
			"http":      3,
			"stdio":     2,
		},
	}

	assert.Equal(t, int64(10), metrics.TotalConnections)
	assert.Equal(t, int64(5), metrics.ActiveConnections)
	assert.Equal(t, int64(100), metrics.MessagesSent)
	assert.Equal(t, int64(150), metrics.MessagesReceived)
	assert.Equal(t, int64(2), metrics.ErrorsTotal)
	assert.Equal(t, int64(1024), metrics.BytesSent)
	assert.Equal(t, int64(2048), metrics.BytesReceived)
	assert.Len(t, metrics.ProtocolStats, 3)
}

// TestConnectionManagerNew tests connection manager creation
func TestConnectionManagerNew(t *testing.T) {
	manager := NewConnectionManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.connections)
	assert.NotNil(t, manager.circuitBreakers)
	assert.NotNil(t, manager.metrics.ProtocolStats)
	assert.NotNil(t, manager.stopCh)
	assert.False(t, manager.running)
}

// TestConnectionManagerStartStop tests start and stop
func TestConnectionManagerStartStop(t *testing.T) {
	manager := NewConnectionManager()

	manager.Start(nil)
	assert.True(t, manager.running)

	manager.Stop(nil)
	assert.False(t, manager.running)
}

// TestCircuitBreakerConstants tests circuit breaker state constants
func TestCircuitBreakerConstants(t *testing.T) {
	assert.Equal(t, CircuitBreakerState("closed"), CircuitBreakerClosed)
	assert.Equal(t, CircuitBreakerState("open"), CircuitBreakerOpen)
	assert.Equal(t, CircuitBreakerState("half_open"), CircuitBreakerHalfOpen)
}

// TestCircuitBreakerConfigStructure tests circuit breaker config
func TestCircuitBreakerConfigStructure(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 5,
		RecoveryTimeout:  30 * time.Second,
		HalfOpenRequests: 3,
	}

	assert.Equal(t, 5, config.FailureThreshold)
	assert.Equal(t, 30*time.Second, config.RecoveryTimeout)
	assert.Equal(t, 3, config.HalfOpenRequests)
}

// TestBaseAdapterNameVersion tests BaseAdapter name and version
func TestBaseAdapterNameVersion(t *testing.T) {
	adapter := &BaseAdapter{
		name:    "test-adapter",
		version: "1.0.0",
	}

	assert.Equal(t, "test-adapter", adapter.Name())
	assert.Equal(t, "1.0.0", adapter.Version())
}

// TestConnectionManagerAddRemove tests connection management
func TestConnectionManagerAddRemove(t *testing.T) {
	manager := NewConnectionManager()

	conn := &AgentConnection{
		ID:       "conn-123",
		Protocol: "websocket",
	}

	manager.Add(conn)
	assert.Equal(t, int64(1), manager.metrics.TotalConnections)
	assert.Equal(t, int64(1), manager.metrics.ActiveConnections)
	assert.Equal(t, int64(1), manager.metrics.ProtocolStats["websocket"])

	manager.Remove("conn-123")
	assert.Equal(t, int64(0), manager.metrics.ActiveConnections)
}

// TestConnectionManagerRemoveNonExistent tests removing non-existent connection
func TestConnectionManagerRemoveNonExistent(t *testing.T) {
	manager := NewConnectionManager()

	assert.NotPanics(t, func() {
		manager.Remove("non-existent")
	})
}

// TestConnectionMetricsZero tests zero metrics
func TestConnectionMetricsZero(t *testing.T) {
	metrics := ConnectionMetrics{}

	assert.Equal(t, int64(0), metrics.TotalConnections)
	assert.Equal(t, int64(0), metrics.ActiveConnections)
	assert.Equal(t, int64(0), metrics.MessagesSent)
	assert.Equal(t, int64(0), metrics.MessagesReceived)
	assert.Equal(t, int64(0), metrics.ErrorsTotal)
	assert.Empty(t, metrics.ProtocolStats)
}

// TestCircuitBreakerConfigZero tests zero config
func TestCircuitBreakerConfigZero(t *testing.T) {
	config := CircuitBreakerConfig{}

	assert.Equal(t, 0, config.FailureThreshold)
	assert.Equal(t, time.Duration(0), config.RecoveryTimeout)
	assert.Equal(t, 0, config.HalfOpenRequests)
}

// TestConnectionManagerMultipleConnections tests multiple connections
func TestConnectionManagerMultipleConnections(t *testing.T) {
	manager := NewConnectionManager()

	for i := 0; i < 5; i++ {
		conn := &AgentConnection{
			ID:       "conn-" + string(rune('1'+i)),
			Protocol: "websocket",
		}
		manager.Add(conn)
	}

	assert.Equal(t, int64(5), manager.metrics.TotalConnections)
	assert.Equal(t, int64(5), manager.metrics.ActiveConnections)
}
