package health

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pryx-core/internal/bus"
)

func setupTestService(t *testing.T) *Service {
	t.Helper()
	b := bus.New()
	return NewService(b)
}

func TestNewService(t *testing.T) {
	svc := setupTestService(t)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.agents)
	assert.NotNil(t, svc.components)
	assert.NotNil(t, svc.alerts)
}

func TestRegisterAgent(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	// Verify agent was registered
	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	assert.Equal(t, "agent-001", agent.AgentID)
	assert.Equal(t, "TestAgent", agent.AgentName)
	assert.Equal(t, "1.0.0", agent.Version)
	assert.Equal(t, HealthStatusHealthy, agent.Status)
}

func TestRegisterAgent_AlreadyExists(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	err = svc.RegisterAgent("agent-001", "TestAgent2", "1.0.1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestDeregisterAgent(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	err = svc.DeregisterAgent("agent-001")
	require.NoError(t, err)

	// Verify agent was deregistered
	_, err = svc.GetAgentHealth("agent-001")
	assert.Error(t, err)
}

func TestDeregisterAgent_NotFound(t *testing.T) {
	svc := setupTestService(t)

	err := svc.DeregisterAgent("non-existent")
	assert.Error(t, err)
}

func TestUpdateHeartbeat(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	// Update heartbeat
	err = svc.UpdateHeartbeat("agent-001")
	require.NoError(t, err)

	// Verify heartbeat was updated
	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	assert.NotZero(t, agent.LastHeartbeat)
}

func TestUpdateHeartbeat_NotFound(t *testing.T) {
	svc := setupTestService(t)

	err := svc.UpdateHeartbeat("non-existent")
	assert.Error(t, err)
}

func TestUpdateComponentHealth(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	component := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeCPU,
		Name:        "CPU Usage",
		Status:      HealthStatusHealthy,
		Metrics:     map[string]interface{}{"usage": 45.5},
	}

	err = svc.UpdateComponentHealth("agent-001", component)
	require.NoError(t, err)

	// Verify component was updated
	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	assert.Len(t, agent.Components, 1)
	assert.Equal(t, "comp-001", agent.Components[0].ComponentID)
}

func TestUpdateComponentHealth_Unhealthy(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	component := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeMemory,
		Name:        "Memory",
		Status:      HealthStatusUnhealthy,
		Metrics:     map[string]interface{}{"usage": 95.0},
		Message:     "High memory usage",
	}

	err = svc.UpdateComponentHealth("agent-001", component)
	require.NoError(t, err)

	// Verify agent status was updated
	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	assert.Equal(t, HealthStatusUnhealthy, agent.Status)
}

func TestPerformHealthCheck(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	req := &HealthCheckRequest{
		RequestID: "check-001",
		AgentID:   "agent-001",
	}

	response, err := svc.PerformHealthCheck(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, "agent-001", response.AgentID)
	assert.Equal(t, HealthStatusHealthy, response.Status)
}

func TestPerformHealthCheck_NotFound(t *testing.T) {
	svc := setupTestService(t)

	req := &HealthCheckRequest{
		RequestID: "check-001",
		AgentID:   "non-existent",
	}

	_, err := svc.PerformHealthCheck(context.Background(), req)
	assert.Error(t, err)
}

func TestGetAgentHealth(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	assert.Equal(t, "agent-001", agent.AgentID)
	assert.Equal(t, HealthStatusHealthy, agent.Status)
}

func TestGetAgentHealth_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetAgentHealth("non-existent")
	assert.Error(t, err)
}

func TestGetAllAgentHealth(t *testing.T) {
	svc := setupTestService(t)

	// Register multiple agents
	for i := 0; i < 3; i++ {
		err := svc.RegisterAgent("agent-"+string(rune('0'+i)), "TestAgent"+string(rune('0'+i)), "1.0.0")
		require.NoError(t, err)
	}

	agents := svc.GetAllAgentHealth()
	assert.Len(t, agents, 3)
}

func TestGetAlerts(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	// Create an unhealthy component to generate an alert
	component := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeCPU,
		Name:        "CPU",
		Status:      HealthStatusUnhealthy,
		Metrics:     map[string]interface{}{"usage": 99.0},
	}

	err = svc.UpdateComponentHealth("agent-001", component)
	require.NoError(t, err)

	// Perform health check to generate alerts
	_, err = svc.PerformHealthCheck(context.Background(), &HealthCheckRequest{
		RequestID: "check-001",
		AgentID:   "agent-001",
	})
	require.NoError(t, err)

	// Get alerts
	alerts := svc.GetAlerts("agent-001", 10)
	assert.Len(t, alerts, 1)
	assert.Equal(t, AlertSeverityCritical, alerts[0].Severity)
}

func TestAcknowledgeAlert(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	// Create an alert
	component := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeMemory,
		Name:        "Memory",
		Status:      HealthStatusUnhealthy,
	}

	err = svc.UpdateComponentHealth("agent-001", component)
	require.NoError(t, err)

	// Perform health check to generate alert
	_, err = svc.PerformHealthCheck(context.Background(), &HealthCheckRequest{
		RequestID: "check-001",
		AgentID:   "agent-001",
	})
	require.NoError(t, err)

	alerts := svc.GetAlerts("agent-001", 10)
	require.Len(t, alerts, 1)

	// Acknowledge the alert
	err = svc.AcknowledgeAlert(alerts[0].AlertID)
	require.NoError(t, err)

	// Verify alert was acknowledged
	alerts = svc.GetAlerts("agent-001", 10)
	assert.Len(t, alerts, 1)
	assert.True(t, alerts[0].Acknowledged)
}

func TestGetAgentCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetAgentCount())

	for i := 0; i < 5; i++ {
		err := svc.RegisterAgent("agent-"+string(rune('0'+i)), "TestAgent", "1.0.0")
		require.NoError(t, err)
	}

	assert.Equal(t, 5, svc.GetAgentCount())
}

func TestGetHealthyCount(t *testing.T) {
	svc := setupTestService(t)

	// Register agents
	for i := 0; i < 3; i++ {
		agentID := fmt.Sprintf("agent-%d", i)
		err := svc.RegisterAgent(agentID, "TestAgent", "1.0.0")
		require.NoError(t, err)
	}

	// Make one unhealthy
	component := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeCPU,
		Name:        "CPU",
		Status:      HealthStatusUnhealthy,
	}
	err := svc.UpdateComponentHealth("agent-0", component)
	require.NoError(t, err)

	assert.Equal(t, 2, svc.GetHealthyCount())
}

func TestGetAlertCount(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	assert.Equal(t, 0, svc.GetAlertCount())

	// Create an unhealthy component
	component := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeDisk,
		Name:        "Disk",
		Status:      HealthStatusUnhealthy,
	}

	err = svc.UpdateComponentHealth("agent-001", component)
	require.NoError(t, err)

	// Perform health check to generate alert
	_, err = svc.PerformHealthCheck(context.Background(), &HealthCheckRequest{
		RequestID: "check-001",
		AgentID:   "agent-001",
	})
	require.NoError(t, err)

	assert.Equal(t, 1, svc.GetAlertCount())
}

func TestPrettyPrint(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	output := agent.PrettyPrint()
	assert.Contains(t, output, "TestAgent")
	assert.Contains(t, output, "agent-001")
	assert.Contains(t, output, "healthy")
}

func TestDegradedStatus(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	// Add a healthy component
	healthy := &HealthComponent{
		ComponentID: "comp-001",
		Type:        ComponentTypeCPU,
		Name:        "CPU",
		Status:      HealthStatusHealthy,
	}
	err = svc.UpdateComponentHealth("agent-001", healthy)
	require.NoError(t, err)

	// Add a degraded component
	degraded := &HealthComponent{
		ComponentID: "comp-002",
		Type:        ComponentTypeMemory,
		Name:        "Memory",
		Status:      HealthStatusDegraded,
	}
	err = svc.UpdateComponentHealth("agent-001", degraded)
	require.NoError(t, err)

	// Verify agent status is degraded
	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)
	assert.Equal(t, HealthStatusDegraded, agent.Status)
}

func TestMultipleComponents(t *testing.T) {
	svc := setupTestService(t)

	err := svc.RegisterAgent("agent-001", "TestAgent", "1.0.0")
	require.NoError(t, err)

	// Add multiple components
	components := []*HealthComponent{
		{ComponentID: "cpu", Type: ComponentTypeCPU, Name: "CPU", Status: HealthStatusHealthy},
		{ComponentID: "memory", Type: ComponentTypeMemory, Name: "Memory", Status: HealthStatusHealthy},
		{ComponentID: "disk", Type: ComponentTypeDisk, Name: "Disk", Status: HealthStatusHealthy},
		{ComponentID: "network", Type: ComponentTypeNetwork, Name: "Network", Status: HealthStatusHealthy},
	}

	for _, comp := range components {
		err = svc.UpdateComponentHealth("agent-001", comp)
		require.NoError(t, err)
	}

	agent, err := svc.GetAgentHealth("agent-001")
	require.NoError(t, err)

	assert.Len(t, agent.Components, 4)
	assert.Equal(t, HealthStatusHealthy, agent.Status)
}
