package trust

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
)

// TrustLevel represents the level of trust between agents
type TrustLevel string

const (
	TrustLevelNone     TrustLevel = "none"
	TrustLevelLow      TrustLevel = "low"
	TrustLevelMedium   TrustLevel = "medium"
	TrustLevelHigh     TrustLevel = "high"
	TrustLevelVerified TrustLevel = "verified"
)

// TrustRelationship represents a trust relationship between two agents
type TrustRelationship struct {
	RelationshipID string                 `json:"relationship_id"`
	SourceAgentID  string                 `json:"source_agent_id"`
	TargetAgentID  string                 `json:"target_agent_id"`
	TrustLevel     TrustLevel             `json:"trust_level"`
	TrustScore     float64                `json:"trust_score"`
	Confidence     float64                `json:"confidence"`
	Factors        []TrustFactor          `json:"factors"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// TrustFactor represents a factor that contributes to trust calculation
type TrustFactor struct {
	FactorID    string  `json:"factor_id"`
	FactorType  string  `json:"factor_type"`
	Weight      float64 `json:"weight"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

// ReputationRecord represents a reputation record for an agent
type ReputationRecord struct {
	RecordID     string                 `json:"record_id"`
	AgentID      string                 `json:"agent_id"`
	EvaluatorID  string                 `json:"evaluator_id"`
	Category     string                 `json:"category"`
	Score        float64                `json:"score"`
	Evidence     []string               `json:"evidence"`
	ReviewedBy   string                 `json:"reviewed_by,omitempty"`
	ReviewStatus string                 `json:"review_status"`
	CreatedAt    time.Time              `json:"created_at"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ReputationScore represents an agent's overall reputation
type ReputationScore struct {
	AgentID        string                 `json:"agent_id"`
	OverallScore   float64                `json:"overall_score"`
	CategoryScores map[string]float64     `json:"category_scores"`
	TrustScore     float64                `json:"trust_score"`
	Reliability    float64                `json:"reliability"`
	Performance    float64                `json:"performance"`
	ReviewsCount   int                    `json:"reviews_count"`
	LastCalculated time.Time              `json:"last_calculated"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// FederationNode represents a node in the federation network
type FederationNode struct {
	NodeID       string                 `json:"node_id"`
	AgentID      string                 `json:"agent_id"`
	AgentName    string                 `json:"agent_name"`
	Endpoint     string                 `json:"endpoint"`
	Status       string                 `json:"status"`
	Capabilities []string               `json:"capabilities"`
	TrustLevel   TrustLevel             `json:"trust_level"`
	LastSeen     time.Time              `json:"last_seen"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// DiscoveryQuery represents a federation discovery query
type DiscoveryQuery struct {
	QueryID     string                 `json:"query_id"`
	RequesterID string                 `json:"requester_id"`
	Criteria    map[string]interface{} `json:"criteria"`
	MaxResults  int                    `json:"max_results"`
	TrustedOnly bool                   `json:"trusted_only"`
	CreatedAt   time.Time              `json:"created_at"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
}

// DiscoveryResult represents the result of a discovery query
type DiscoveryResult struct {
	QueryID    string            `json:"query_id"`
	Nodes      []*FederationNode `json:"nodes"`
	TotalFound int               `json:"total_found"`
	ReturnedAt time.Time         `json:"returned_at"`
}

// TrustPolicy represents a trust policy for federation decisions
type TrustPolicy struct {
	PolicyID      string      `json:"policy_id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	MinTrustLevel TrustLevel  `json:"min_trust_level"`
	RequiredCaps  []string    `json:"required_caps"`
	AllowedCats   []string    `json:"allowed_categories"`
	Rules         []TrustRule `json:"rules"`
	Priority      int         `json:"priority"`
	Enabled       bool        `json:"enabled"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// TrustRule represents a single trust rule
type TrustRule struct {
	RuleID    string  `json:"rule_id"`
	RuleType  string  `json:"rule_type"`
	Condition string  `json:"condition"`
	Action    string  `json:"action"`
	Weight    float64 `json:"weight"`
}

// Service manages trust, reputation, and federation
type Service struct {
	mu             sync.RWMutex
	bus            *bus.Bus
	relationships  map[string]*TrustRelationship
	reputations    map[string]*ReputationScore
	records        map[string][]*ReputationRecord
	nodes          map[string]*FederationNode
	queries        map[string]*DiscoveryQuery
	policies       map[string]*TrustPolicy
	propagationLog []TrustEvent
}

// NewService creates a new trust management service
func NewService(b *bus.Bus) *Service {
	return &Service{
		bus:            b,
		relationships:  make(map[string]*TrustRelationship),
		reputations:    make(map[string]*ReputationScore),
		records:        make(map[string][]*ReputationRecord),
		nodes:          make(map[string]*FederationNode),
		queries:        make(map[string]*DiscoveryQuery),
		policies:       make(map[string]*TrustPolicy),
		propagationLog: make([]TrustEvent, 0),
	}
}

// EstablishTrust establishes a trust relationship between two agents
func (s *Service) EstablishTrust(ctx context.Context, sourceID, targetID string, level TrustLevel, factors []TrustFactor) (*TrustRelationship, error) {
	if sourceID == targetID {
		return nil, fmt.Errorf("cannot establish trust with self")
	}

	rel := &TrustRelationship{
		RelationshipID: uuid.New().String(),
		SourceAgentID:  sourceID,
		TargetAgentID:  targetID,
		TrustLevel:     level,
		TrustScore:     s.calculateTrustScore(level, factors),
		Confidence:     s.calculateConfidence(factors),
		Factors:        factors,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Metadata:       make(map[string]interface{}),
	}

	s.mu.Lock()
	s.relationships[rel.RelationshipID] = rel
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("trust.relationship.established", "", map[string]interface{}{
		"relationship_id": rel.RelationshipID,
		"source_agent":    sourceID,
		"target_agent":    targetID,
		"trust_level":     level,
		"trust_score":     rel.TrustScore,
	}))

	return rel, nil
}

// GetTrustRelationship retrieves a trust relationship
func (s *Service) GetTrustRelationship(sourceID, targetID string) (*TrustRelationship, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, rel := range s.relationships {
		if rel.SourceAgentID == sourceID && rel.TargetAgentID == targetID {
			return rel, nil
		}
	}

	return nil, fmt.Errorf("trust relationship not found between %s and %s", sourceID, targetID)
}

// UpdateTrustLevel updates the trust level for a relationship
func (s *Service) UpdateTrustLevel(relationshipID string, level TrustLevel, factors []TrustFactor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rel, exists := s.relationships[relationshipID]
	if !exists {
		return fmt.Errorf("relationship not found: %s", relationshipID)
	}

	rel.TrustLevel = level
	rel.TrustScore = s.calculateTrustScore(level, factors)
	rel.Factors = factors
	rel.UpdatedAt = time.Now().UTC()

	// Publish event
	s.bus.Publish(bus.NewEvent("trust.relationship.updated", "", map[string]interface{}{
		"relationship_id": relationshipID,
		"trust_level":     level,
		"trust_score":     rel.TrustScore,
	}))

	return nil
}

// GetTrustedAgents retrieves all agents trusted by a given agent
func (s *Service) GetTrustedAgents(agentID string, minLevel TrustLevel) []*TrustRelationship {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var trusted []*TrustRelationship
	for _, rel := range s.relationships {
		if rel.SourceAgentID == agentID && s.meetsTrustThreshold(rel.TrustLevel, minLevel) {
			trusted = append(trusted, rel)
		}
	}

	return trusted
}

// calculateTrustScore calculates a trust score based on level and factors
func (s *Service) calculateTrustScore(level TrustLevel, factors []TrustFactor) float64 {
	baseScore := map[TrustLevel]float64{
		TrustLevelNone:     0.0,
		TrustLevelLow:      0.25,
		TrustLevelMedium:   0.50,
		TrustLevelHigh:     0.75,
		TrustLevelVerified: 1.0,
	}[level]

	factorScore := 0.0
	factorWeight := 0.0

	for _, f := range factors {
		factorScore += f.Score * f.Weight
		factorWeight += f.Weight
	}

	if factorWeight > 0 {
		return math.Min(1.0, baseScore*0.7+factorScore/factorWeight*0.3)
	}

	return baseScore
}

// calculateConfidence calculates confidence score based on factors
func (s *Service) calculateConfidence(factors []TrustFactor) float64 {
	if len(factors) == 0 {
		return 0.5
	}

	total := 0.0
	for _, f := range factors {
		total += f.Weight
	}

	return math.Min(1.0, total/float64(len(factors)))
}

// meetsTrustThreshold checks if a trust level meets the minimum threshold
func (s *Service) meetsTrustThreshold(actual, minimum TrustLevel) bool {
	levels := []TrustLevel{
		TrustLevelNone,
		TrustLevelLow,
		TrustLevelMedium,
		TrustLevelHigh,
		TrustLevelVerified,
	}

	actualIdx := -1
	minIdx := -1

	for i, l := range levels {
		if l == actual {
			actualIdx = i
		}
		if l == minimum {
			minIdx = i
		}
	}

	return actualIdx >= minIdx && actualIdx >= 0 && minIdx >= 0
}

// RecordReputation records a reputation entry for an agent
func (s *Service) RecordReputation(ctx context.Context, record *ReputationRecord) error {
	record.RecordID = uuid.New().String()
	record.CreatedAt = time.Now().UTC()

	s.mu.Lock()
	s.records[record.AgentID] = append(s.records[record.AgentID], record)
	s.mu.Unlock()

	// Recalculate reputation score
	s.recalculateReputation(record.AgentID)

	// Publish event
	s.bus.Publish(bus.NewEvent("trust.reputation.recorded", "", map[string]interface{}{
		"record_id":    record.RecordID,
		"agent_id":     record.AgentID,
		"category":     record.Category,
		"score":        record.Score,
		"evaluator_id": record.EvaluatorID,
	}))

	return nil
}

// GetReputationScore retrieves the reputation score for an agent
func (s *Service) GetReputationScore(agentID string) (*ReputationScore, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	score, exists := s.reputations[agentID]
	if !exists {
		return nil, fmt.Errorf("reputation not found for agent: %s", agentID)
	}

	return score, nil
}

// recalculateReputation recalculates the reputation score for an agent
func (s *Service) recalculateReputation(agentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	records := s.records[agentID]
	if len(records) == 0 {
		return
	}

	categoryScores := make(map[string]float64)
	categoryCounts := make(map[string]int)
	overallSum := 0.0
	reviewsCount := 0

	for _, rec := range records {
		if rec.ReviewStatus == "approved" || rec.ReviewStatus == "pending" {
			categoryScores[rec.Category] += rec.Score
			categoryCounts[rec.Category]++
			overallSum += rec.Score
			reviewsCount++
		}
	}

	// Calculate averages for each category
	for cat, sum := range categoryScores {
		if count := categoryCounts[cat]; count > 0 {
			categoryScores[cat] = sum / float64(count)
		}
	}

	score := &ReputationScore{
		AgentID:        agentID,
		OverallScore:   overallSum / float64(reviewsCount),
		CategoryScores: categoryScores,
		TrustScore:     categoryScores["trust"],
		Reliability:    categoryScores["reliability"],
		Performance:    categoryScores["performance"],
		ReviewsCount:   reviewsCount,
		LastCalculated: time.Now().UTC(),
		Metadata:       make(map[string]interface{}),
	}

	s.reputations[agentID] = score
}

// RegisterNode registers a node in the federation network
func (s *Service) RegisterNode(node *FederationNode) error {
	node.NodeID = uuid.New().String()
	node.LastSeen = time.Now().UTC()

	s.mu.Lock()
	s.nodes[node.AgentID] = node
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("trust.node.registered", "", map[string]interface{}{
		"node_id":     node.NodeID,
		"agent_id":    node.AgentID,
		"agent_name":  node.AgentName,
		"endpoint":    node.Endpoint,
		"trust_level": node.TrustLevel,
	}))

	return nil
}

// UpdateNodeHeartbeat updates the last seen time for a node
func (s *Service) UpdateNodeHeartbeat(agentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, exists := s.nodes[agentID]
	if !exists {
		return fmt.Errorf("node not found: %s", agentID)
	}

	node.LastSeen = time.Now().UTC()

	return nil
}

// DiscoverNodes discovers federation nodes matching criteria
func (s *Service) DiscoverNodes(ctx context.Context, query *DiscoveryQuery) (*DiscoveryResult, error) {
	query.QueryID = uuid.New().String()
	query.CreatedAt = time.Now().UTC()

	var matching []*FederationNode
	s.mu.RLock()
	for _, node := range s.nodes {
		if s.matchesQuery(node, query) {
			if len(matching) < query.MaxResults {
				matching = append(matching, node)
			}
		}
	}
	s.mu.RUnlock()

	result := &DiscoveryResult{
		QueryID:    query.QueryID,
		Nodes:      matching,
		TotalFound: len(matching),
		ReturnedAt: time.Now().UTC(),
	}

	// Store query
	s.mu.Lock()
	s.queries[query.QueryID] = query
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("trust.discovery.completed", "", map[string]interface{}{
		"query_id":     query.QueryID,
		"requester_id": query.RequesterID,
		"found_count":  len(matching),
	}))

	return result, nil
}

// matchesQuery checks if a node matches the discovery query
func (s *Service) matchesQuery(node *FederationNode, query *DiscoveryQuery) bool {
	if query.TrustedOnly && node.TrustLevel != TrustLevelHigh && node.TrustLevel != TrustLevelVerified {
		return false
	}

	for key, value := range query.Criteria {
		switch key {
		case "status":
			if node.Status != value.(string) {
				return false
			}
		case "capability":
			found := false
			for _, cap := range node.Capabilities {
				if cap == value.(string) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case "trust_level":
			if node.TrustLevel != TrustLevel(value.(string)) {
				return false
			}
		}
	}

	return true
}

// PropagateTrust propagates trust to connected agents
func (s *Service) PropagateTrust(relationshipID string) error {
	s.mu.RLock()
	rel, exists := s.relationships[relationshipID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("relationship not found: %s", relationshipID)
	}

	// Find agents that trust the source agent
	var toUpdate []*TrustRelationship
	s.mu.RLock()
	for _, r := range s.relationships {
		if r.TargetAgentID == rel.SourceAgentID && r.RelationshipID != relationshipID {
			toUpdate = append(toUpdate, r)
		}
	}
	s.mu.RUnlock()

	// Update their trust scores based on the new relationship
	for _, updateRel := range toUpdate {
		newScore := (updateRel.TrustScore + rel.TrustScore) / 2
		s.mu.Lock()
		updateRel.TrustScore = math.Min(1.0, newScore)
		updateRel.UpdatedAt = time.Now().UTC()
		s.mu.Unlock()

		// Log propagation
		event := TrustEvent{
			EventID:          uuid.New().String(),
			FromRelationship: relationshipID,
			ToRelationship:   updateRel.RelationshipID,
			ScoreChange:      newScore - updateRel.TrustScore,
			PropagatedAt:     time.Now().UTC(),
		}
		s.mu.Lock()
		s.propagationLog = append(s.propagationLog, event)
		s.mu.Unlock()
	}

	// Publish event
	s.bus.Publish(bus.NewEvent("trust.propagated", "", map[string]interface{}{
		"relationship_id": relationshipID,
		"updated_count":   len(toUpdate),
	}))

	return nil
}

// TrustEvent represents a trust propagation event
type TrustEvent struct {
	EventID          string    `json:"event_id"`
	FromRelationship string    `json:"from_relationship"`
	ToRelationship   string    `json:"to_relationship"`
	ScoreChange      float64   `json:"score_change"`
	PropagatedAt     time.Time `json:"propagated_at"`
}

// AddPolicy adds a trust policy
func (s *Service) AddPolicy(policy *TrustPolicy) error {
	policy.PolicyID = uuid.New().String()
	policy.CreatedAt = time.Now().UTC()
	policy.UpdatedAt = time.Now().UTC()

	s.mu.Lock()
	s.policies[policy.PolicyID] = policy
	s.mu.Unlock()

	return nil
}

// EvaluatePolicy evaluates if an agent meets trust policy requirements
func (s *Service) EvaluatePolicy(agentID string, policyID string) (bool, error) {
	s.mu.RLock()
	policy, exists := s.policies[policyID]
	s.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("policy not found: %s", policyID)
	}

	rep, err := s.GetReputationScore(agentID)
	if err != nil {
		return false, err
	}

	// Determine trust level from score
	currentLevel := s.scoreToTrustLevel(rep.TrustScore)

	// Check minimum trust level
	if !s.meetsTrustThreshold(currentLevel, policy.MinTrustLevel) {
		return false, nil
	}

	// Check required capabilities
	for _, cap := range policy.RequiredCaps {
		found := false
		s.mu.RLock()
		for _, node := range s.nodes {
			if node.AgentID == agentID {
				for _, nodeCap := range node.Capabilities {
					if nodeCap == cap {
						found = true
						break
					}
				}
				break
			}
		}
		s.mu.RUnlock()

		if !found {
			return false, nil
		}
	}

	return true, nil
}

// scoreToTrustLevel converts a numeric score to a trust level
func (s *Service) scoreToTrustLevel(score float64) TrustLevel {
	switch {
	case score >= 0.9:
		return TrustLevelVerified
	case score >= 0.7:
		return TrustLevelHigh
	case score >= 0.5:
		return TrustLevelMedium
	case score >= 0.3:
		return TrustLevelLow
	default:
		return TrustLevelNone
	}
}

// GetStats returns trust service statistics
func (s *Service) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"relationships_count": len(s.relationships),
		"reputations_count":   len(s.reputations),
		"nodes_count":         len(s.nodes),
		"policies_count":      len(s.policies),
		"pending_queries":     len(s.queries),
		"propagation_events":  len(s.propagationLog),
	}
}

// PrettyPrint prints trust relationship info
func (t *TrustRelationship) PrettyPrint() string {
	return fmt.Sprintf("TrustRelationship{Source: %s -> Target: %s, Level: %s, Score: %.2f}",
		t.SourceAgentID, t.TargetAgentID, t.TrustLevel, t.TrustScore)
}

// PrettyPrint prints reputation score info
func (r *ReputationScore) PrettyPrint() string {
	return fmt.Sprintf("ReputationScore{Agent: %s, Overall: %.2f, Trust: %.2f, Reviews: %d}",
		r.AgentID, r.OverallScore, r.TrustScore, r.ReviewsCount)
}

// MarshalJSON for TrustRelationship
func (t *TrustRelationship) MarshalJSON() ([]byte, error) {
	type Alias TrustRelationship
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

// MarshalJSON for ReputationScore
func (r *ReputationScore) MarshalJSON() ([]byte, error) {
	type Alias ReputationScore
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}
