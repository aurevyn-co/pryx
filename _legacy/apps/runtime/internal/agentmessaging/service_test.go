package agentmessaging

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pryx-core/internal/bus"
	"pryx-core/internal/message"
)

func setupTestService(t *testing.T) (*Service, *bus.Bus, *message.Service) {
	t.Helper()
	b := bus.New()
	msgSvc := message.NewService(b, "test-agent")
	svc := NewService(b, msgSvc)
	return svc, b, msgSvc
}

func TestNewService(t *testing.T) {
	svc, _, _ := setupTestService(t)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.bus)
	assert.NotNil(t, svc.msgService)
	assert.NotNil(t, svc.sessions)
	assert.NotNil(t, svc.conversations)
	assert.NotNil(t, svc.history)
}

func TestConnectSession(t *testing.T) {
	svc, _, _ := setupTestService(t)

	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)
	require.NotNil(t, session)

	assert.NotEmpty(t, session.SessionID)
	assert.Equal(t, "agent-001", session.RemoteAgentID)
	assert.Equal(t, "TestAgent", session.RemoteAgentName)
	assert.Equal(t, "http://localhost:8080", session.RemoteEndpoint)
	assert.Equal(t, SessionStatusConnected, session.Status)
	assert.NotZero(t, session.CreatedAt)
	assert.NotZero(t, session.LastActivity)
}

func TestConnectSession_StoresSession(t *testing.T) {
	svc, _, _ := setupTestService(t)

	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)

	retrieved, err := svc.GetSession(session.SessionID)
	require.NoError(t, err)
	assert.Equal(t, session.SessionID, retrieved.SessionID)
}

func TestDisconnectSession(t *testing.T) {
	svc, _, _ := setupTestService(t)

	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)

	err = svc.DisconnectSession(session.SessionID)
	require.NoError(t, err)

	retrieved, err := svc.GetSession(session.SessionID)
	require.NoError(t, err)
	assert.Equal(t, SessionStatusDisconnected, retrieved.Status)
}

func TestDisconnectSession_NotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	err := svc.DisconnectSession("non-existent-session")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

func TestSendMessage(t *testing.T) {
	svc, _, msgSvc := setupTestService(t)

	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)

	// Register a handler for test messages
	msgSvc.RegisterHandler(message.MessageTypeRequest, func(ctx context.Context, msg *message.Message) (*message.Response, error) {
		return message.NewResponse(msg.CorrelationID, 200, "OK", msg.Payload), nil
	})

	receipt, err := svc.SendMessage(context.Background(), session.SessionID, message.MessageTypeRequest, map[string]interface{}{
		"content": "Hello",
	})
	require.NoError(t, err)
	require.NotNil(t, receipt)

	assert.NotEmpty(t, receipt.ReceiptID)
	assert.NotEmpty(t, receipt.MessageID)
	assert.Equal(t, "agent-001", receipt.ToAgent)
	assert.Equal(t, "delivered", receipt.Status)
}

func TestSendMessage_SessionNotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	receipt, err := svc.SendMessage(context.Background(), "non-existent-session", message.MessageTypeRequest, nil)
	assert.Error(t, err)
	assert.Nil(t, receipt)
	assert.Contains(t, err.Error(), "session not found")
}

func TestSendMessage_SessionNotConnected(t *testing.T) {
	svc, _, _ := setupTestService(t)

	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)

	// Disconnect the session
	err = svc.DisconnectSession(session.SessionID)
	require.NoError(t, err)

	receipt, err := svc.SendMessage(context.Background(), session.SessionID, message.MessageTypeRequest, nil)
	assert.Error(t, err)
	assert.Nil(t, receipt)
	assert.Contains(t, err.Error(), "session not connected")
}

func TestSendDirectMessage(t *testing.T) {
	svc, _, msgSvc := setupTestService(t)
	_ = msgSvc

	// Register a handler

	// Register a handler
	msgSvc.RegisterHandler(message.MessageTypeRequest, func(ctx context.Context, msg *message.Message) (*message.Response, error) {
		return message.NewResponse(msg.CorrelationID, 200, "OK", msg.Payload), nil
	})

	receipt, err := svc.SendDirectMessage(context.Background(), "agent-001", message.MessageTypeRequest, map[string]interface{}{
		"message": "direct message",
	})
	require.NoError(t, err)
	require.NotNil(t, receipt)

	assert.NotEmpty(t, receipt.MessageID)
	assert.Equal(t, "agent-001", receipt.ToAgent)
	assert.Equal(t, "delivered", receipt.Status)
}

func TestReplyToMessage(t *testing.T) {
	svc, _, msgSvc := setupTestService(t)
	_ = msgSvc

	// Create a conversation first

	// Create a conversation first
	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	// Add a message to the conversation
	testMsg := message.NewMessage(message.MessageTypeRequest, "agent-001", "agent-002", map[string]interface{}{
		"content": "original message",
	})
	err = svc.AddMessageToConversation(context.Background(), conv.ConversationID, "agent-001", testMsg)
	require.NoError(t, err)

	// Register handler for responses
	msgSvc.RegisterHandler(message.MessageTypeResponse, func(ctx context.Context, msg *message.Message) (*message.Response, error) {
		return message.NewResponse(msg.CorrelationID, 200, "OK", msg.Payload), nil
	})

	// Reply to the message
	receipt, err := svc.ReplyToMessage(context.Background(), testMsg.ID, map[string]interface{}{
		"content": "reply message",
	})
	require.NoError(t, err)
	require.NotNil(t, receipt)
}

func TestStartConversation(t *testing.T) {
	svc, _, _ := setupTestService(t)

	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002", "agent-003"}, "Test Subject")
	require.NoError(t, err)
	require.NotNil(t, conv)

	assert.NotEmpty(t, conv.ConversationID)
	assert.Equal(t, "Test Subject", conv.Subject)
	assert.Len(t, conv.Participants, 3)
	assert.Equal(t, ConversationStatusActive, conv.Status)
	assert.NotZero(t, conv.CreatedAt)
	assert.NotZero(t, conv.UpdatedAt)
}

func TestStartConversation_StoresConversation(t *testing.T) {
	svc, _, _ := setupTestService(t)

	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	retrieved, err := svc.GetConversation(conv.ConversationID)
	require.NoError(t, err)
	assert.Equal(t, conv.ConversationID, retrieved.ConversationID)
}

func TestGetConversationsByParticipant(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create multiple conversations
	_, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Conv 1")
	require.NoError(t, err)

	_, err = svc.StartConversation(context.Background(), []string{"agent-001", "agent-003"}, "Conv 2")
	require.NoError(t, err)

	_, err = svc.StartConversation(context.Background(), []string{"agent-002", "agent-003"}, "Conv 3")
	require.NoError(t, err)

	// Get conversations for agent-001
	convs := svc.GetConversationsByParticipant("agent-001")
	assert.Len(t, convs, 2)

	// Get conversations for agent-002
	convs = svc.GetConversationsByParticipant("agent-002")
	assert.Len(t, convs, 2)

	// Get conversations for agent-003
	convs = svc.GetConversationsByParticipant("agent-003")
	assert.Len(t, convs, 2)
}

func TestAddMessageToConversation(t *testing.T) {
	svc, _, _ := setupTestService(t)

	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	testMsg := message.NewMessage(message.MessageTypeRequest, "agent-001", "agent-002", map[string]interface{}{
		"content": "Hello, agent-002!",
	})

	err = svc.AddMessageToConversation(context.Background(), conv.ConversationID, "agent-001", testMsg)
	require.NoError(t, err)

	// Verify message was added
	retrieved, err := svc.GetConversation(conv.ConversationID)
	require.NoError(t, err)
	assert.Len(t, retrieved.Messages, 1)
	assert.Equal(t, testMsg.ID, retrieved.Messages[0].ID)
}

func TestAddMessageToConversation_NotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	testMsg := message.NewMessage(message.MessageTypeRequest, "agent-001", "agent-002", map[string]interface{}{
		"content": "Hello",
	})

	err := svc.AddMessageToConversation(context.Background(), "non-existent-conversation", "agent-001", testMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conversation not found")
}

func TestArchiveConversation(t *testing.T) {
	svc, _, _ := setupTestService(t)

	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	err = svc.ArchiveConversation(conv.ConversationID)
	require.NoError(t, err)

	retrieved, err := svc.GetConversation(conv.ConversationID)
	require.NoError(t, err)
	assert.Equal(t, ConversationStatusArchived, retrieved.Status)
}

func TestGetMessageHistory(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create a conversation with messages
	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	// Add multiple messages
	for i := 0; i < 5; i++ {
		testMsg := message.NewMessage(message.MessageTypeRequest, "agent-001", "agent-002", map[string]interface{}{
			"content": "Message",
		})
		err = svc.AddMessageToConversation(context.Background(), conv.ConversationID, "agent-001", testMsg)
		require.NoError(t, err)
	}

	// Get message history for agent-001
	history := svc.GetMessageHistory("agent-001", 10)
	assert.Len(t, history, 5)

	// Get with limit
	history = svc.GetMessageHistory("agent-001", 3)
	assert.Len(t, history, 3)
}

func TestGetActiveSessions(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create multiple sessions
	_, err := svc.ConnectSession(context.Background(), "agent-001", "Agent1", "http://localhost:8081")
	require.NoError(t, err)

	_, err = svc.ConnectSession(context.Background(), "agent-002", "Agent2", "http://localhost:8082")
	require.NoError(t, err)

	_, err = svc.ConnectSession(context.Background(), "agent-003", "Agent3", "http://localhost:8083")
	require.NoError(t, err)

	sessions := svc.GetActiveSessions()
	assert.Len(t, sessions, 3)

	// Disconnect one session
	_, err = svc.ConnectSession(context.Background(), "agent-004", "Agent4", "http://localhost:8084")
	require.NoError(t, err)
	err = svc.DisconnectSession(sessions[0].SessionID)
	require.NoError(t, err)

	sessions = svc.GetActiveSessions()
	assert.Len(t, sessions, 3)
}

func TestGetSessionsByAgent(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create sessions with the same agent
	_, err := svc.ConnectSession(context.Background(), "agent-001", "Agent1", "http://localhost:8081")
	require.NoError(t, err)

	_, err = svc.ConnectSession(context.Background(), "agent-001", "Agent1", "http://localhost:8082")
	require.NoError(t, err)

	_, err = svc.ConnectSession(context.Background(), "agent-002", "Agent2", "http://localhost:8083")
	require.NoError(t, err)

	sessions := svc.GetSessionsByAgent("agent-001")
	assert.Len(t, sessions, 2)

	sessions = svc.GetSessionsByAgent("agent-002")
	assert.Len(t, sessions, 1)
}

func TestMarkMessageRead(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create a conversation with a message
	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	testMsg := message.NewMessage(message.MessageTypeRequest, "agent-001", "agent-002", map[string]interface{}{
		"content": "Hello",
	})
	err = svc.AddMessageToConversation(context.Background(), conv.ConversationID, "agent-001", testMsg)
	require.NoError(t, err)

	// Mark message as read
	err = svc.MarkMessageRead(testMsg.ID)
	require.NoError(t, err)
}

func TestMarkMessageRead_NotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	err := svc.MarkMessageRead("non-existent-message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message not found")
}

func TestPrettyPrint(t *testing.T) {
	svc, _, _ := setupTestService(t)

	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)

	output := session.PrettyPrint()
	assert.Contains(t, output, "Session")
	assert.Contains(t, output, "agent-001")
	assert.Contains(t, output, "TestAgent")
	assert.Contains(t, output, "connected")

	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	output = conv.PrettyPrint()
	assert.Contains(t, output, "Conversation")
	assert.Contains(t, output, "Test Subject")
	assert.Contains(t, output, "active")
}

func TestGetActiveSessionCount(t *testing.T) {
	svc, _, _ := setupTestService(t)

	assert.Equal(t, 0, svc.GetActiveSessionCount())

	_, err := svc.ConnectSession(context.Background(), "agent-001", "Agent1", "http://localhost:8081")
	require.NoError(t, err)
	assert.Equal(t, 1, svc.GetActiveSessionCount())

	_, err = svc.ConnectSession(context.Background(), "agent-002", "Agent2", "http://localhost:8082")
	require.NoError(t, err)
	assert.Equal(t, 2, svc.GetActiveSessionCount())
}

func TestGetConversationCount(t *testing.T) {
	svc, _, _ := setupTestService(t)

	assert.Equal(t, 0, svc.GetConversationCount())

	_, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Conv 1")
	require.NoError(t, err)
	assert.Equal(t, 1, svc.GetConversationCount())

	_, err = svc.StartConversation(context.Background(), []string{"agent-001", "agent-003"}, "Conv 2")
	require.NoError(t, err)
	assert.Equal(t, 2, svc.GetConversationCount())
}

func TestGetConversation_NotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	_, err := svc.GetConversation("non-existent-conversation")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conversation not found")
}

func TestGetSession_NotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	_, err := svc.GetSession("non-existent-session")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

func TestMessageHistory_Tags(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create a conversation with a message
	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	testMsg := message.NewMessage(message.MessageTypeRequest, "agent-001", "agent-002", map[string]interface{}{
		"content": "Hello",
	})
	err = svc.AddMessageToConversation(context.Background(), conv.ConversationID, "agent-001", testMsg)
	require.NoError(t, err)

	// Check that the history entry has the conversation ID as a tag
	history := svc.GetMessageHistory("agent-001", 10)
	require.Len(t, history, 1)
	require.Len(t, history[0].Tags, 1)
	assert.Equal(t, conv.ConversationID, history[0].Tags[0])
}

func TestSessionConcurrency(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create a session
	session, err := svc.ConnectSession(context.Background(), "agent-001", "TestAgent", "http://localhost:8080")
	require.NoError(t, err)

	// Concurrently disconnect and get the session
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = svc.GetSession(session.SessionID)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent operations")
		}
	}
}

func TestConversationConcurrency(t *testing.T) {
	svc, _, _ := setupTestService(t)

	// Create a conversation
	conv, err := svc.StartConversation(context.Background(), []string{"agent-001", "agent-002"}, "Test Subject")
	require.NoError(t, err)

	// Concurrently get the conversation
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = svc.GetConversation(conv.ConversationID)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent operations")
		}
	}
}
