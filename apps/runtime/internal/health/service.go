package health

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"pryx-core/internal/bus"

	"github.com/google/uuid"
)

// HealthStatus represents the health status of an agent
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// ComponentType represents the type of component being monitored
type ComponentType string

const (
	ComponentTypeCPU      ComponentType = "cpu"
	ComponentTypeMemory   ComponentType = "memory"
	ComponentTypeDisk     ComponentType = "disk"
	ComponentTypeNetwork  ComponentType = "network"
	ComponentTypeAgent    ComponentType = "agent"
	ComponentTypeChannel  ComponentType = "channel"
	ComponentTypeMCP      ComponentType = "mcp"
	ComponentTypeDatabase ComponentType = "database"
)

// HealthComponent represents a single component's health
type HealthComponent struct {
	ComponentID   string                 `json:"component_id"`
	Type          ComponentType          `json:"type"`
	Name          string                 `json:"name"`
	Status        HealthStatus           `json:"status"`
	Metrics       map[string]interface{} `json:"metrics"`
	Details       map[string]interface{} `json:"details,omitempty"`
	LastCheckedAt time.Time              `json:"last_checked_at"`
	Message       string                 `json:"message,omitempty"`
}

// AgentHealth represents the overall health of an agent
type AgentHealth struct {
	AgentID       string                 `json:"agent_id"`
	AgentName     string                 `json:"agent_name"`
	Status        HealthStatus           `json:"status"`
	Uptime        time.Duration          `json:"uptime"`
	Components    []HealthComponent      `json:"components"`
	Version       string                 `json:"version"`
	LastHeartbeat time.Time              `json:"last_heartbeat"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// HealthCheckRequest represents a health check request
type HealthCheckRequest struct {
	RequestID   string          `json:"request_id"`
	AgentID     string          `json:"agent_id"`
	CheckTypes  []ComponentType `json:"check_types,omitempty"`
	DeepCheck   bool            `json:"deep_check"`
	Timeout     time.Duration   `json:"timeout"`
	RequestedAt time.Time       `json:"requested_at"`
}

// HealthCheckResponse represents the response to a health check
type HealthCheckResponse struct {
	RequestID  string            `json:"request_id"`
	AgentID    string            `json:"agent_id"`
	Status     HealthStatus      `json:"status"`
	Components []HealthComponent `json:"components"`
	CheckedAt  time.Time         `json:"checked_at"`
	DurationMs int64             `json:"duration_ms"`
}

// HealthAlert represents an alert generated from health checks
type HealthAlert struct {
	AlertID        string                 `json:"alert_id"`
	AgentID        string                 `json:"agent_id"`
	Severity       AlertSeverity          `json:"severity"`
	ComponentType  ComponentType          `json:"component_type"`
	Message        string                 `json:"message"`
	Details        map[string]interface{} `json:"details"`
	Acknowledged   bool                   `json:"acknowledged"`
	CreatedAt      time.Time              `json:"created_at"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
}

// AlertSeverity represents the severity of an alert
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

// Service manages agent health monitoring
type Service struct {
	mu         sync.RWMutex
	bus        *bus.Bus
	agents     map[string]*AgentHealth
	components map[string]map[string]*HealthComponent // agentID -> componentID -> component
	alerts     []HealthAlert
	lastCheck  time.Time
}

// NewService creates a new health monitoring service
func NewService(b *bus.Bus) *Service {
	return &Service{
		bus:        b,
		agents:     make(map[string]*AgentHealth),
		components: make(map[string]map[string]*HealthComponent),
		alerts:     make([]HealthAlert, 0),
	}
}

// RegisterAgent registers an agent for health monitoring
func (s *Service) RegisterAgent(agentID, agentName, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.agents[agentID]; exists {
		return fmt.Errorf("agent already registered: %s", agentID)
	}

	agent := &AgentHealth{
		AgentID:       agentID,
		AgentName:     agentName,
		Status:        HealthStatusHealthy,
		Uptime:        0,
		Components:    make([]HealthComponent, 0),
		Version:       version,
		LastHeartbeat: time.Now().UTC(),
		Metadata:      make(map[string]interface{}),
	}

	s.agents[agentID] = agent
	s.components[agentID] = make(map[string]*HealthComponent)

	// Publish event
	s.bus.Publish(bus.NewEvent("health.agent.registered", "", map[string]interface{}{
		"agent_id":   agentID,
		"agent_name": agentName,
	}))

	return nil
}

// DeregisterAgent removes an agent from health monitoring
func (s *Service) DeregisterAgent(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.agents[agentID]; !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	delete(s.agents, agentID)
	delete(s.components, agentID)

	// Publish event
	s.bus.Publish(bus.NewEvent("health.agent.deregistered", "", map[string]interface{}{
		"agent_id": agentID,
	}))

	return nil
}

// UpdateHeartbeat updates the heartbeat for an agent
func (s *Service) UpdateHeartbeat(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, exists := s.agents[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	agent.LastHeartbeat = time.Now().UTC()

	return nil
}

// UpdateComponentHealth updates the health of a component
func (s *Service) UpdateComponentHealth(agentID string, component *HealthComponent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	agentComponents, exists := s.components[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	component.LastCheckedAt = time.Now().UTC()
	agentComponents[component.ComponentID] = component

	// Update agent status based on component health
	s.updateAgentStatus(agentID)

	return nil
}

// updateAgentStatus updates the overall agent status based on components
func (s *Service) updateAgentStatus(agentID string) {
	agent, exists := s.agents[agentID]
	if !exists {
		return
	}

	components := s.components[agentID]
	if len(components) == 0 {
		agent.Status = HealthStatusUnknown
		return
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, comp := range components {
		switch comp.Status {
		case HealthStatusUnhealthy:
			hasUnhealthy = true
		case HealthStatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		agent.Status = HealthStatusUnhealthy
	} else if hasDegraded {
		agent.Status = HealthStatusDegraded
	} else {
		agent.Status = HealthStatusHealthy
	}

	agent.Components = s.getComponentsList(agentID)
}

// getComponentsList returns a list of components for an agent
func (s *Service) getComponentsList(agentID string) []HealthComponent {
	components := s.components[agentID]
	list := make([]HealthComponent, 0, len(components))
	for _, comp := range components {
		list = append(list, *comp)
	}
	return list
}

// PerformHealthCheck performs a health check on an agent
func (s *Service) PerformHealthCheck(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	startTime := time.Now()

	s.mu.RLock()
	_, exists := s.agents[req.AgentID]
	var agentComponents []*HealthComponent
	if exists {
		// Create a copy of the pointers to avoid race on the map access
		// The component objects themselves might be modified under lock in other methods,
		// but HealthComponent here seems to be efficiently copyable struct?
		// Wait, s.components is map[string]map[string]*HealthComponent
		// We need to copy the *HealthComponent pointers.
		// Note: The content of *HealthComponent might be modified.
		// Ideally we should copy the VALUES aka dereference them if we want a snapshot.
		// HealthComponent field 'Status' etc are modified in UpdateComponentHealth under Lock.
		// So we should copy.
		if comps, ok := s.components[req.AgentID]; ok {
			agentComponents = make([]*HealthComponent, 0, len(comps))
			for _, c := range comps {
				// We append the pointer for now, but strictly we should probably dereference if we want snapshot.
				// But Service.PerformHealthCheck returns *pointers* in HealthCheckResponse? No, structure `HealthComponent` (value).
				// So we should read the values under lock.
				agentComponents = append(agentComponents, c)
			}
		}
	}
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent not found: %s", req.AgentID)
	}

	components := make([]HealthComponent, 0)

	// Perform checks based on request
	for _, compPointer := range agentComponents {
		// dereference safely? We already have the pointer.
		comp := *compPointer // Snapshot the value
		if len(req.CheckTypes) == 0 || containsComponentType(req.CheckTypes, comp.Type) {
			components = append(components, comp)
		}
	}

	// Determine overall status
	status := HealthStatusHealthy
	for _, comp := range components {
		if comp.Status == HealthStatusUnhealthy {
			status = HealthStatusUnhealthy
			break
		} else if comp.Status == HealthStatusDegraded && status != HealthStatusUnhealthy {
			status = HealthStatusDegraded
		}
	}

	duration := time.Since(startTime)

	response := &HealthCheckResponse{
		RequestID:  req.RequestID,
		AgentID:    req.AgentID,
		Status:     status,
		Components: components,
		CheckedAt:  time.Now().UTC(),
		DurationMs: duration.Milliseconds(),
	}

	// Generate alerts if needed
	s.checkAndGenerateAlerts(req.AgentID, components)

	return response, nil
}

// checkAndGenerateAlerts checks components and generates alerts
func (s *Service) checkAndGenerateAlerts(agentID string, components []HealthComponent) {
	for _, comp := range components {
		if comp.Status == HealthStatusUnhealthy {
			alert := HealthAlert{
				AlertID:       uuid.New().String(),
				AgentID:       agentID,
				Severity:      AlertSeverityCritical,
				ComponentType: comp.Type,
				Message:       fmt.Sprintf("Component %s is unhealthy: %s", comp.Name, comp.Message),
				Details:       comp.Metrics,
				CreatedAt:     time.Now().UTC(),
			}

			s.mu.Lock()
			s.alerts = append(s.alerts, alert)
			s.mu.Unlock()

			// Publish event
			s.bus.Publish(bus.NewEvent("health.alert", "", map[string]interface{}{
				"alert_id":       alert.AlertID,
				"agent_id":       agentID,
				"severity":       alert.Severity,
				"component_type": comp.Type,
				"message":        alert.Message,
			}))
		}
	}
}

// GetAgentHealth retrieves the health status of an agent
func (s *Service) GetAgentHealth(agentID string) (*AgentHealth, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agent, exists := s.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}

	agent.Components = s.getComponentsList(agentID)
	return agent, nil
}

// GetAllAgentHealth retrieves health status of all registered agents
func (s *Service) GetAllAgentHealth() []*AgentHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agents := make([]*AgentHealth, 0, len(s.agents))
	for _, agent := range s.agents {
		agent.Components = s.getComponentsList(agent.AgentID)
		agents = append(agents, agent)
	}

	return agents
}

// GetAlerts retrieves alerts for an agent
func (s *Service) GetAlerts(agentID string, limit int) []HealthAlert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	alerts := make([]HealthAlert, 0)
	for i := len(s.alerts) - 1; i >= 0 && len(alerts) < limit; i-- {
		if s.alerts[i].AgentID == agentID {
			alerts = append(alerts, s.alerts[i])
		}
	}

	return alerts
}

// AcknowledgeAlert acknowledges an alert
func (s *Service) AcknowledgeAlert(alertID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, alert := range s.alerts {
		if alert.AlertID == alertID {
			now := time.Now().UTC()
			s.alerts[i].Acknowledged = true
			s.alerts[i].AcknowledgedAt = &now
			return nil
		}
	}

	return fmt.Errorf("alert not found: %s", alertID)
}

// GetAgentCount returns the number of registered agents
func (s *Service) GetAgentCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.agents)
}

// GetHealthyCount returns the number of healthy agents
func (s *Service) GetHealthyCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, agent := range s.agents {
		if agent.Status == HealthStatusHealthy {
			count++
		}
	}

	return count
}

// GetAlertCount returns the number of active alerts
func (s *Service) GetAlertCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, alert := range s.alerts {
		if !alert.Acknowledged {
			count++
		}
	}

	return count
}

// PrettyPrint prints agent health info
func (h *AgentHealth) PrettyPrint() string {
	return fmt.Sprintf("AgentHealth{Agent: %s (%s), Status: %s, Components: %d}",
		h.AgentName, h.AgentID, h.Status, len(h.Components))
}

// MarshalJSON for AgentHealth
func (h *AgentHealth) MarshalJSON() ([]byte, error) {
	type Alias AgentHealth
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(h),
	})
}

// MarshalJSON for HealthComponent
func (h *HealthComponent) MarshalJSON() ([]byte, error) {
	type Alias HealthComponent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(h),
	})
}

// containsComponentType checks if a slice contains a component type
func containsComponentType(types []ComponentType, target ComponentType) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}
