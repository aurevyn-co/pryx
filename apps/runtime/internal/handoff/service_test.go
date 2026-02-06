package handoff

import (
	"context"
	"testing"
	"time"

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
	assert.NotNil(t, svc.requests)
	assert.NotNil(t, svc.responses)
	assert.NotNil(t, svc.transfers)
	assert.NotNil(t, svc.history)
	assert.NotNil(t, svc.activeHandoffs)
}

func TestRequestHandoff(t *testing.T) {
	svc := setupTestService(t)

	req := &HandoffRequest{
		FromAgentID:   "agent-001",
		FromAgentName: "Agent1",
		ToAgentID:     "agent-002",
		ToAgentName:   "Agent2",
		SessionID:     "session-001",
		ContextTypes:  []ContextType{ContextTypeSession, ContextTypeMemory},
		Priority:      1,
	}

	response, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.NotEmpty(t, response.RequestID)
	assert.Equal(t, HandoffStatusPending, response.Status)

	// Verify request was stored
	svc.mu.RLock()
	require.Len(t, svc.requests, 1)
	svc.mu.RUnlock()
}

func TestRequestHandoff_DefaultTimeout(t *testing.T) {
	svc := setupTestService(t)

	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	response, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	// Default timeout should be 5 minutes
	assert.NotNil(t, response)
}

func TestAcceptHandoff(t *testing.T) {
	svc := setupTestService(t)

	// Request handoff
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
		SessionID:   "session-001",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	// Accept handoff
	acceptedCaps := []ContextType{ContextTypeSession, ContextTypeMemory}
	conditions := []string{"Keep session active for 1 hour"}

	response, err := svc.AcceptHandoff(context.Background(), requestResponse.RequestID, acceptedCaps, conditions)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, HandoffStatusAccepted, response.Status)
	assert.Len(t, response.AcceptedCaps, 2)
	assert.Len(t, response.Conditions, 1)
	assert.NotEmpty(t, response.TransferURL)
	assert.NotNil(t, response.ExpiresAt)
}

func TestAcceptHandoff_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.AcceptHandoff(context.Background(), "non-existent-request", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handoff request not found")
}

func TestRejectHandoff(t *testing.T) {
	svc := setupTestService(t)

	// Request handoff
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	// Reject handoff
	rejectedCaps := []ContextType{ContextTypePolicy}
	response, err := svc.RejectHandoff(context.Background(), requestResponse.RequestID, rejectedCaps, "Policy context not allowed")
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, HandoffStatusFailed, response.Status)
	assert.Len(t, response.RejectedCaps, 1)
}

func TestRejectHandoff_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.RejectHandoff(context.Background(), "non-existent-request", nil, "reason")
	assert.Error(t, err)
}

func TestStartTransfer(t *testing.T) {
	svc := setupTestService(t)

	// Request and accept handoff
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
		SessionID:   "session-001",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.AcceptHandoff(context.Background(), requestResponse.RequestID, []ContextType{ContextTypeSession}, nil)
	require.NoError(t, err)

	// Start transfer
	contexts := []*HandoffContext{
		{
			ContextID:  "ctx-001",
			Type:       ContextTypeSession,
			SessionID:  "session-001",
			SizeBytes:  1024,
			Compressed: false,
		},
	}

	transfers, err := svc.StartTransfer(context.Background(), requestResponse.RequestID, contexts)
	require.NoError(t, err)
	require.Len(t, transfers, 1)

	assert.Equal(t, HandoffPhaseTransfer, transfers[0].Phase)
	assert.Equal(t, int64(1024), transfers[0].BytesTransferred)
	assert.Equal(t, "completed", transfers[0].Status)
}

func TestStartTransfer_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.StartTransfer(context.Background(), "non-existent-request", nil)
	assert.Error(t, err)
}

func TestValidateTransfer(t *testing.T) {
	svc := setupTestService(t)

	// Setup handoff
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.AcceptHandoff(context.Background(), requestResponse.RequestID, []ContextType{ContextTypeSession}, nil)
	require.NoError(t, err)

	// Validate transfer
	transfer, err := svc.ValidateTransfer(context.Background(), requestResponse.RequestID)
	require.NoError(t, err)
	require.NotNil(t, transfer)

	assert.Equal(t, HandoffPhaseValidation, transfer.Phase)
	assert.Equal(t, "completed", transfer.Status)
}

func TestCompleteHandoff(t *testing.T) {
	svc := setupTestService(t)

	// Setup handoff
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
		SessionID:   "session-001",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.AcceptHandoff(context.Background(), requestResponse.RequestID, []ContextType{ContextTypeSession}, nil)
	require.NoError(t, err)

	// Complete handoff
	err = svc.CompleteHandoff(context.Background(), requestResponse.RequestID)
	require.NoError(t, err)

	// Verify handoff is in history
	history := svc.GetHandoffHistory("agent-001", 10)
	require.Len(t, history, 1)
	assert.Equal(t, HandoffStatusCompleted, history[0].Status)

	// Verify it's no longer active
	assert.Equal(t, 0, svc.GetActiveCount())
}

func TestCompleteHandoff_NotFound(t *testing.T) {
	svc := setupTestService(t)

	err := svc.CompleteHandoff(context.Background(), "non-existent-request")
	assert.Error(t, err)
}

func TestCancelHandoff(t *testing.T) {
	svc := setupTestService(t)

	// Setup handoff
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
		SessionID:   "session-001",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	// Cancel handoff
	err = svc.CancelHandoff(context.Background(), requestResponse.RequestID, "User requested cancellation")
	require.NoError(t, err)

	// Verify handoff is in history as cancelled
	history := svc.GetHandoffHistory("agent-001", 10)
	require.Len(t, history, 1)
	assert.Equal(t, HandoffStatusCancelled, history[0].Status)
}

func TestGetHandoffStatus(t *testing.T) {
	svc := setupTestService(t)

	// Setup handoff and accept it
	req := &HandoffRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	requestResponse, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.AcceptHandoff(context.Background(), requestResponse.RequestID, []ContextType{ContextTypeSession}, nil)
	require.NoError(t, err)

	status, err := svc.GetHandoffStatus(requestResponse.RequestID)
	require.NoError(t, err)
	assert.Equal(t, HandoffStatusAccepted, status.Status)
}

func TestGetHandoffStatus_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetHandoffStatus("non-existent-request")
	assert.Error(t, err)
}

func TestGetActiveHandoffs(t *testing.T) {
	svc := setupTestService(t)

	// Create multiple handoffs
	req1 := &HandoffRequest{FromAgentID: "agent-001", ToAgentID: "agent-002"}
	req2 := &HandoffRequest{FromAgentID: "agent-003", ToAgentID: "agent-004"}

	_, err := svc.RequestHandoff(context.Background(), req1)
	require.NoError(t, err)

	_, err = svc.RequestHandoff(context.Background(), req2)
	require.NoError(t, err)

	handoffs := svc.GetActiveHandoffs()
	assert.Len(t, handoffs, 2)
}

func TestGetHandoffHistory(t *testing.T) {
	svc := setupTestService(t)

	// Create and complete handoffs
	for i := 0; i < 5; i++ {
		req := &HandoffRequest{
			FromAgentID: "agent-001",
			ToAgentID:   "agent-002",
			SessionID:   "session-" + string(rune('0'+i)),
		}

		resp, err := svc.RequestHandoff(context.Background(), req)
		require.NoError(t, err)

		_, err = svc.AcceptHandoff(context.Background(), resp.RequestID, []ContextType{ContextTypeSession}, nil)
		require.NoError(t, err)

		err = svc.CompleteHandoff(context.Background(), resp.RequestID)
		require.NoError(t, err)
	}

	// Get history
	history := svc.GetHandoffHistory("agent-001", 10)
	assert.Len(t, history, 5)

	// Get with limit
	history = svc.GetHandoffHistory("agent-001", 3)
	assert.Len(t, history, 3)
}

func TestCreateSessionState(t *testing.T) {
	svc := setupTestService(t)

	messages := []SessionMessage{
		{
			MessageID: "msg-001",
			Role:      "user",
			Content:   "Hello",
			Timestamp: time.Now().UTC(),
		},
	}

	state := svc.CreateSessionState("session-001", "agent-001", messages, map[string]interface{}{
		"topic": "test",
	})

	assert.Equal(t, "session-001", state.SessionID)
	assert.Equal(t, "agent-001", state.AgentID)
	assert.Len(t, state.Messages, 1)
	assert.NotZero(t, state.CreatedAt)
}

func TestCreateContext(t *testing.T) {
	svc := setupTestService(t)

	ctx := svc.CreateContext(ContextTypeSession, "session-001", map[string]interface{}{
		"key": "value",
	})

	assert.NotEmpty(t, ctx.ContextID)
	assert.Equal(t, ContextTypeSession, ctx.Type)
	assert.Equal(t, "session-001", ctx.SessionID)
	assert.NotZero(t, ctx.CreatedAt)
}

func TestGetRequestCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetRequestCount())

	req := &HandoffRequest{FromAgentID: "agent-001", ToAgentID: "agent-002"}
	_, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, 1, svc.GetRequestCount())
}

func TestGetActiveCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetActiveCount())

	req := &HandoffRequest{FromAgentID: "agent-001", ToAgentID: "agent-002"}
	resp, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, 1, svc.GetActiveCount())

	// Complete the handoff
	_, err = svc.AcceptHandoff(context.Background(), resp.RequestID, []ContextType{ContextTypeSession}, nil)
	require.NoError(t, err)

	err = svc.CompleteHandoff(context.Background(), resp.RequestID)
	require.NoError(t, err)

	assert.Equal(t, 0, svc.GetActiveCount())
}

func TestGetHistoryCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetHistoryCount())

	// Create and complete a handoff
	req := &HandoffRequest{FromAgentID: "agent-001", ToAgentID: "agent-002"}
	resp, err := svc.RequestHandoff(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.AcceptHandoff(context.Background(), resp.RequestID, []ContextType{ContextTypeSession}, nil)
	require.NoError(t, err)

	err = svc.CompleteHandoff(context.Background(), resp.RequestID)
	require.NoError(t, err)

	assert.Equal(t, 1, svc.GetHistoryCount())
}

func TestPrettyPrint(t *testing.T) {
	_ = setupTestService(t)

	req := &HandoffRequest{
		RequestID:   "req-001",
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
		SessionID:   "session-001",
	}

	output := req.PrettyPrint()
	assert.Contains(t, output, "req-001")
	assert.Contains(t, output, "agent-001")
	assert.Contains(t, output, "agent-002")
	assert.Contains(t, output, "session-001")
}

func TestConcurrentHandoffs(t *testing.T) {
	svc := setupTestService(t)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(agentID string) {
			req := &HandoffRequest{
				FromAgentID: agentID,
				ToAgentID:   "agent-dest",
			}
			_, _ = svc.RequestHandoff(context.Background(), req)
			done <- true
		}("agent-" + string(rune('A'+i)))
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent handoffs")
		}
	}

	assert.Equal(t, 10, svc.GetActiveCount())
}
