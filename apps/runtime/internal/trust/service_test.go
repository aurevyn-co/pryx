package trust

import (
	"context"
	"sync"
	"testing"
	"time"

	"pryx-core/internal/bus"
)

// helperService creates a trust service for testing
func helperService(t *testing.T) (*Service, *bus.Bus) {
	b := bus.New()
	s := NewService(b)
	return s, b
}

// TestEstablishTrust tests establishing a trust relationship
func TestEstablishTrust(t *testing.T) {
	s, _ := helperService(t)

	factors := []TrustFactor{
		{FactorID: "1", FactorType: "history", Weight: 0.5, Score: 0.8, Description: "Positive history"},
		{FactorID: "2", FactorType: "recommendation", Weight: 0.5, Score: 0.9, Description: "Trusted recommendation"},
	}

	rel, err := s.EstablishTrust(context.Background(), "agent-a", "agent-b", TrustLevelHigh, factors)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rel.RelationshipID == "" {
		t.Error("expected relationship ID to be set")
	}
	if rel.SourceAgentID != "agent-a" {
		t.Errorf("expected source agent 'agent-a', got '%s'", rel.SourceAgentID)
	}
	if rel.TargetAgentID != "agent-b" {
		t.Errorf("expected target agent 'agent-b', got '%s'", rel.TargetAgentID)
	}
	if rel.TrustLevel != TrustLevelHigh {
		t.Errorf("expected trust level High, got %s", rel.TrustLevel)
	}
	if rel.TrustScore <= 0 {
		t.Error("expected trust score to be positive")
	}
}

// TestEstablishTrustSelf tests that self-trust is rejected
func TestEstablishTrustSelf(t *testing.T) {
	s, _ := helperService(t)

	_, err := s.EstablishTrust(context.Background(), "agent-a", "agent-a", TrustLevelHigh, nil)
	if err == nil {
		t.Error("expected error when establishing trust with self")
	}
}

// TestGetTrustRelationship tests retrieving a trust relationship
func TestGetTrustRelationship(t *testing.T) {
	s, _ := helperService(t)

	// First establish a trust relationship
	_, err := s.EstablishTrust(context.Background(), "agent-a", "agent-b", TrustLevelMedium, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Then retrieve it
	rel, err := s.GetTrustRelationship("agent-a", "agent-b")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rel.SourceAgentID != "agent-a" || rel.TargetAgentID != "agent-b" {
		t.Error("retrieved wrong relationship")
	}
}

// TestGetTrustRelationshipNotFound tests error when relationship not found
func TestGetTrustRelationshipNotFound(t *testing.T) {
	s, _ := helperService(t)

	_, err := s.GetTrustRelationship("agent-x", "agent-y")
	if err == nil {
		t.Error("expected error when relationship not found")
	}
}

// TestUpdateTrustLevel tests updating trust level
func TestUpdateTrustLevel(t *testing.T) {
	s, _ := helperService(t)

	// Establish initial trust
	rel, err := s.EstablishTrust(context.Background(), "agent-a", "agent-b", TrustLevelLow, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Update trust level
	factors := []TrustFactor{
		{FactorID: "1", FactorType: "verified", Weight: 1.0, Score: 1.0, Description: "Identity verified"},
	}

	err = s.UpdateTrustLevel(rel.RelationshipID, TrustLevelVerified, factors)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify update
	rel, err = s.GetTrustRelationship("agent-a", "agent-b")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rel.TrustLevel != TrustLevelVerified {
		t.Errorf("expected trust level Verified, got %s", rel.TrustLevel)
	}
}

// TestGetTrustedAgents tests retrieving trusted agents
func TestGetTrustedAgents(t *testing.T) {
	s, _ := helperService(t)

	// Establish multiple trust relationships
	s.EstablishTrust(context.Background(), "agent-a", "agent-b", TrustLevelLow, nil)
	s.EstablishTrust(context.Background(), "agent-a", "agent-c", TrustLevelHigh, nil)
	s.EstablishTrust(context.Background(), "agent-a", "agent-d", TrustLevelMedium, nil)

	// Get high-trust agents
	trusted := s.GetTrustedAgents("agent-a", TrustLevelHigh)
	if len(trusted) != 1 {
		t.Errorf("expected 1 high-trust agent, got %d", len(trusted))
	}

	// Get medium-trust agents
	trusted = s.GetTrustedAgents("agent-a", TrustLevelMedium)
	if len(trusted) != 2 {
		t.Errorf("expected 2 medium-trust agents, got %d", len(trusted))
	}
}

// TestCalculateTrustScore tests trust score calculation
func TestCalculateTrustScore(t *testing.T) {
	s, _ := helperService(t)

	testCases := []struct {
		level    TrustLevel
		minScore float64
		maxScore float64
	}{
		{TrustLevelNone, 0.0, 0.35},
		{TrustLevelLow, 0.25, 0.5},
		{TrustLevelMedium, 0.45, 0.7},
		{TrustLevelHigh, 0.65, 0.9},
		{TrustLevelVerified, 0.85, 1.0},
	}

	for _, tc := range testCases {
		factors := []TrustFactor{
			{FactorID: "1", FactorType: "test", Weight: 0.5, Score: 0.8, Description: "test"},
			{FactorID: "2", FactorType: "test", Weight: 0.5, Score: 0.8, Description: "test"},
		}

		score := s.calculateTrustScore(tc.level, factors)
		if score < tc.minScore || score > tc.maxScore {
			t.Errorf("for level %s: score %.2f not in range [%.2f, %.2f]",
				tc.level, score, tc.minScore, tc.maxScore)
		}
	}
}

// TestMeetsTrustThreshold tests trust threshold checking
func TestMeetsTrustThreshold(t *testing.T) {
	s, _ := helperService(t)

	testCases := []struct {
		actual   TrustLevel
		minimum  TrustLevel
		expected bool
	}{
		{TrustLevelLow, TrustLevelNone, true},
		{TrustLevelMedium, TrustLevelLow, true},
		{TrustLevelHigh, TrustLevelMedium, true},
		{TrustLevelVerified, TrustLevelHigh, true},
		{TrustLevelNone, TrustLevelLow, false},
		{TrustLevelLow, TrustLevelMedium, false},
		{TrustLevelMedium, TrustLevelHigh, false},
		{TrustLevelHigh, TrustLevelVerified, false},
	}

	for _, tc := range testCases {
		result := s.meetsTrustThreshold(tc.actual, tc.minimum)
		if result != tc.expected {
			t.Errorf("meetsTrustThreshold(%s, %s) = %v, expected %v",
				tc.actual, tc.minimum, result, tc.expected)
		}
	}
}

// TestRecordReputation tests recording reputation entries
func TestRecordReputation(t *testing.T) {
	s, _ := helperService(t)

	record := &ReputationRecord{
		AgentID:      "agent-b",
		EvaluatorID:  "agent-a",
		Category:     "reliability",
		Score:        0.85,
		Evidence:     []string{"completed task successfully"},
		ReviewStatus: "approved",
	}

	err := s.RecordReputation(context.Background(), record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check reputation was recorded
	s.mu.RLock()
	records := s.records["agent-b"]
	s.mu.RUnlock()

	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
}

// TestGetReputationScore tests retrieving reputation scores
func TestGetReputationScore(t *testing.T) {
	s, _ := helperService(t)

	// Record some reputation entries
	s.RecordReputation(context.Background(), &ReputationRecord{
		AgentID:      "agent-b",
		EvaluatorID:  "agent-a",
		Category:     "trust",
		Score:        0.8,
		ReviewStatus: "approved",
	})

	s.RecordReputation(context.Background(), &ReputationRecord{
		AgentID:      "agent-b",
		EvaluatorID:  "agent-c",
		Category:     "reliability",
		Score:        0.9,
		ReviewStatus: "approved",
	})

	score, err := s.GetReputationScore("agent-b")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if score.AgentID != "agent-b" {
		t.Errorf("expected agent ID 'agent-b', got '%s'", score.AgentID)
	}
	if score.OverallScore <= 0 {
		t.Error("expected positive overall score")
	}
	if score.ReviewsCount != 2 {
		t.Errorf("expected 2 reviews, got %d", score.ReviewsCount)
	}
}

// TestGetReputationScoreNotFound tests error when reputation not found
func TestGetReputationScoreNotFound(t *testing.T) {
	s, _ := helperService(t)

	_, err := s.GetReputationScore("unknown-agent")
	if err == nil {
		t.Error("expected error when reputation not found")
	}
}

// TestRegisterNode tests registering a federation node
func TestRegisterNode(t *testing.T) {
	s, _ := helperService(t)

	node := &FederationNode{
		AgentID:      "agent-b",
		AgentName:    "Agent B",
		Endpoint:     "https://agent-b.example.com",
		Status:       "online",
		Capabilities: []string{"text-generation", "code-analysis"},
		TrustLevel:   TrustLevelHigh,
	}

	err := s.RegisterNode(node)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	s.mu.RLock()
	n := s.nodes["agent-b"]
	s.mu.RUnlock()

	if n.AgentID != "agent-b" {
		t.Errorf("expected agent ID 'agent-b', got '%s'", n.AgentID)
	}
	if n.NodeID == "" {
		t.Error("expected node ID to be set")
	}
}

// TestUpdateNodeHeartbeat tests updating node heartbeat
func TestUpdateNodeHeartbeat(t *testing.T) {
	s, _ := helperService(t)

	s.RegisterNode(&FederationNode{
		AgentID:   "agent-b",
		AgentName: "Agent B",
	})

	initialLastSeen := s.nodes["agent-b"].LastSeen

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	err := s.UpdateNodeHeartbeat("agent-b")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !s.nodes["agent-b"].LastSeen.After(initialLastSeen) {
		t.Error("expected last seen to be updated")
	}
}

// TestUpdateNodeHeartbeatNotFound tests error when node not found
func TestUpdateNodeHeartbeatNotFound(t *testing.T) {
	s, _ := helperService(t)

	err := s.UpdateNodeHeartbeat("unknown-agent")
	if err == nil {
		t.Error("expected error when node not found")
	}
}

// TestDiscoverNodes tests discovering federation nodes
func TestDiscoverNodes(t *testing.T) {
	s, _ := helperService(t)

	// Register some nodes
	s.RegisterNode(&FederationNode{
		AgentID:      "agent-b",
		AgentName:    "Agent B",
		Capabilities: []string{"text-generation"},
		TrustLevel:   TrustLevelHigh,
	})

	s.RegisterNode(&FederationNode{
		AgentID:      "agent-c",
		AgentName:    "Agent C",
		Capabilities: []string{"code-analysis"},
		TrustLevel:   TrustLevelMedium,
	})

	query := &DiscoveryQuery{
		RequesterID: "agent-a",
		Criteria: map[string]interface{}{
			"capability": "text-generation",
		},
		MaxResults: 10,
	}

	result, err := s.DiscoverNodes(context.Background(), query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(result.Nodes))
	}

	if result.Nodes[0].AgentID != "agent-b" {
		t.Errorf("expected agent-b, got %s", result.Nodes[0].AgentID)
	}
}

// TestDiscoverNodesTrustedOnly tests discovery with trusted-only filter
func TestDiscoverNodesTrustedOnly(t *testing.T) {
	s, _ := helperService(t)

	// Register nodes with different trust levels
	s.RegisterNode(&FederationNode{
		AgentID:    "agent-b",
		AgentName:  "Agent B",
		TrustLevel: TrustLevelHigh,
	})

	s.RegisterNode(&FederationNode{
		AgentID:    "agent-c",
		AgentName:  "Agent C",
		TrustLevel: TrustLevelLow,
	})

	query := &DiscoveryQuery{
		RequesterID: "agent-a",
		Criteria:    map[string]interface{}{},
		MaxResults:  10,
		TrustedOnly: true,
	}

	result, err := s.DiscoverNodes(context.Background(), query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Nodes) != 1 {
		t.Errorf("expected 1 trusted node, got %d", len(result.Nodes))
	}
}

// TestPropagateTrust tests trust propagation
func TestPropagateTrust(t *testing.T) {
	s, _ := helperService(t)

	// Establish initial trust: A -> B
	relAB, _ := s.EstablishTrust(context.Background(), "agent-a", "agent-b", TrustLevelHigh, nil)

	// Establish trust: C -> A (so C trusts A, and trust should propagate to B)
	s.EstablishTrust(context.Background(), "agent-c", "agent-a", TrustLevelMedium, nil)

	err := s.PropagateTrust(relAB.RelationshipID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check propagation log
	s.mu.RLock()
	logLen := len(s.propagationLog)
	s.mu.RUnlock()

	if logLen == 0 {
		t.Error("expected propagation events in log")
	}
}

// TestAddPolicy tests adding trust policies
func TestAddPolicy(t *testing.T) {
	s, _ := helperService(t)

	policy := &TrustPolicy{
		Name:          "High Security Policy",
		Description:   "Requires verified trust and specific capabilities",
		MinTrustLevel: TrustLevelHigh,
		RequiredCaps:  []string{"security-audit"},
		Priority:      100,
		Enabled:       true,
	}

	err := s.AddPolicy(policy)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if policy.PolicyID == "" {
		t.Error("expected policy ID to be set")
	}
}

// TestEvaluatePolicy tests policy evaluation
func TestEvaluatePolicy(t *testing.T) {
	s, _ := helperService(t)

	// Add a policy
	policy := &TrustPolicy{
		Name:          "Test Policy",
		MinTrustLevel: TrustLevelMedium,
		RequiredCaps:  []string{"text-generation"},
		Priority:      1,
		Enabled:       true,
	}
	s.AddPolicy(policy)

	// Register a node with required capabilities
	s.RegisterNode(&FederationNode{
		AgentID:      "agent-b",
		AgentName:    "Agent B",
		Capabilities: []string{"text-generation"},
		TrustLevel:   TrustLevelMedium,
	})

	// Record reputation that meets the medium trust level
	s.RecordReputation(context.Background(), &ReputationRecord{
		AgentID:      "agent-b",
		EvaluatorID:  "agent-a",
		Category:     "trust",
		Score:        0.65,
		ReviewStatus: "approved",
	})

	pass, err := s.EvaluatePolicy("agent-b", policy.PolicyID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !pass {
		t.Error("expected policy evaluation to pass")
	}
}

// TestEvaluatePolicyFailTrustLevel tests policy evaluation failure due to trust level
func TestEvaluatePolicyFailTrustLevel(t *testing.T) {
	s, _ := helperService(t)

	// Add a policy requiring high trust
	policy := &TrustPolicy{
		Name:          "High Trust Policy",
		MinTrustLevel: TrustLevelVerified,
		Priority:      1,
		Enabled:       true,
	}
	s.AddPolicy(policy)

	// Register node with low trust
	s.RegisterNode(&FederationNode{
		AgentID:    "agent-b",
		AgentName:  "Agent B",
		TrustLevel: TrustLevelLow,
	})

	// Record low reputation
	s.RecordReputation(context.Background(), &ReputationRecord{
		AgentID:      "agent-b",
		EvaluatorID:  "agent-a",
		Category:     "trust",
		Score:        0.3,
		ReviewStatus: "approved",
	})

	pass, err := s.EvaluatePolicy("agent-b", policy.PolicyID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pass {
		t.Error("expected policy evaluation to fail due to low trust")
	}
}

// TestScoreToTrustLevel tests score to trust level conversion
func TestScoreToTrustLevel(t *testing.T) {
	s, _ := helperService(t)

	testCases := []struct {
		score    float64
		expected TrustLevel
	}{
		{0.95, TrustLevelVerified},
		{0.85, TrustLevelHigh},
		{0.75, TrustLevelHigh},
		{0.65, TrustLevelMedium},
		{0.50, TrustLevelMedium},
		{0.40, TrustLevelLow},
		{0.30, TrustLevelLow},
		{0.20, TrustLevelNone},
	}

	for _, tc := range testCases {
		level := s.scoreToTrustLevel(tc.score)
		if level != tc.expected {
			t.Errorf("score %.2f: expected %s, got %s",
				tc.score, tc.expected, level)
		}
	}
}

// TestGetStats tests getting service statistics
func TestGetStats(t *testing.T) {
	s, _ := helperService(t)

	// Add some data
	s.EstablishTrust(context.Background(), "a", "b", TrustLevelHigh, nil)
	s.RegisterNode(&FederationNode{AgentID: "b"})
	s.AddPolicy(&TrustPolicy{Name: "test"})

	stats := s.GetStats()

	if stats["relationships_count"].(int) != 1 {
		t.Errorf("expected 1 relationship, got %v", stats["relationships_count"])
	}
	if stats["nodes_count"].(int) != 1 {
		t.Errorf("expected 1 node, got %v", stats["nodes_count"])
	}
	if stats["policies_count"].(int) != 1 {
		t.Errorf("expected 1 policy, got %v", stats["policies_count"])
	}
}

// TestPrettyPrint tests PrettyPrint methods
func TestPrettyPrint(t *testing.T) {
	rel := &TrustRelationship{
		SourceAgentID: "agent-a",
		TargetAgentID: "agent-b",
		TrustLevel:    TrustLevelHigh,
		TrustScore:    0.85,
	}

	output := rel.PrettyPrint()
	if output == "" {
		t.Error("expected non-empty output")
	}

	rep := &ReputationScore{
		AgentID:      "agent-b",
		OverallScore: 0.9,
		TrustScore:   0.85,
		ReviewsCount: 5,
	}

	output = rep.PrettyPrint()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

// TestConcurrentAccess tests thread-safe access to the service
func TestConcurrentAccess(t *testing.T) {
	s, _ := helperService(t)

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrent trust establishment
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			s.EstablishTrust(context.Background(),
				"agent-a",
				"agent-b",
				TrustLevelHigh,
				nil)
		}(i)
	}

	// Concurrent reputation recording
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			s.RecordReputation(context.Background(), &ReputationRecord{
				AgentID:      "agent-b",
				EvaluatorID:  "agent-x",
				Category:     "reliability",
				Score:        0.8,
				ReviewStatus: "approved",
			})
		}(i)
	}

	// Concurrent node registration
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			s.RegisterNode(&FederationNode{
				AgentID:   "agent-c",
				AgentName: "Agent C",
			})
		}(i)
	}

	wg.Wait()

	// Verify data integrity
	if s.nodes["agent-c"] == nil {
		t.Error("expected node to be registered")
	}
}

// TestMarshalJSON tests JSON marshaling
func TestMarshalJSON(t *testing.T) {
	rel := &TrustRelationship{
		RelationshipID: "rel-1",
		SourceAgentID:  "agent-a",
		TargetAgentID:  "agent-b",
		TrustLevel:     TrustLevelHigh,
		TrustScore:     0.85,
		CreatedAt:      time.Now(),
	}

	data, err := rel.MarshalJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}

	rep := &ReputationScore{
		AgentID:        "agent-b",
		OverallScore:   0.9,
		CategoryScores: map[string]float64{"trust": 0.85},
		ReviewsCount:   5,
	}

	data, err = rep.MarshalJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
}

// TestEventPublishing tests that events are published correctly
func TestEventPublishing(t *testing.T) {
	b := bus.New()
	s := NewService(b)

	// Subscribe to trust events
	events := make([]bus.Event, 0)
	var mu sync.Mutex
	eventCh, closer := b.Subscribe()
	defer closer()

	go func() {
		for event := range eventCh {
			mu.Lock()
			events = append(events, event)
			mu.Unlock()
		}
	}()

	// Perform an action that should publish an event
	s.EstablishTrust(context.Background(), "agent-a", "agent-b", TrustLevelHigh, nil)

	// Wait for event
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	eventCount := len(events)
	mu.Unlock()
	if eventCount == 0 {
		t.Error("expected at least one event to be published")
	}
}
