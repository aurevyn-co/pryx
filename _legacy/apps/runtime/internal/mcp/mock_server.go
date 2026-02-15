package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type MockServer struct {
	mu           sync.RWMutex
	tools        []Tool
	initialized  atomic.Bool
	callCount    map[string]int
	lastCallArgs map[string]map[string]interface{}

	InitializeFunc func(ctx context.Context, req RPCRequest) RPCResponse
	ListToolsFunc  func(ctx context.Context) ([]Tool, error)
	CallToolFunc   func(ctx context.Context, name string, args map[string]interface{}) (ToolResult, error)
}

func NewMockServer() *MockServer {
	m := &MockServer{
		tools: []Tool{
			{
				Name:        "echo",
				Title:       "Echo Tool",
				Description: "Echoes back the input",
				InputSchema: json.RawMessage(`{"type":"object","properties":{"message":{"type":"string"}},"required":["message"]}`),
			},
			{
				Name:        "add",
				Title:       "Add Numbers",
				Description: "Adds two numbers",
				InputSchema: json.RawMessage(`{"type":"object","properties":{"a":{"type":"number"},"b":{"type":"number"}},"required":["a","b"]}`),
			},
		},
		callCount:    make(map[string]int),
		lastCallArgs: make(map[string]map[string]interface{}),
	}

	m.InitializeFunc = m.defaultInitialize
	m.ListToolsFunc = m.defaultListTools
	m.CallToolFunc = m.defaultCallTool

	return m
}

func (m *MockServer) HandleRequest(ctx context.Context, req RPCRequest) RPCResponse {
	switch req.Method {
	case "initialize":
		return m.handleInitialize(ctx, req)
	case "tools/list":
		return m.handleListTools(ctx, req)
	case "tools/call":
		return m.handleCallTool(ctx, req)
	case "ping":
		return m.handlePing(ctx, req)
	default:
		return RPCResponse{
			JSONRPC: "2.0",
			ID:      mustMarshalID(req.ID),
			Error:   &RPCError{Code: -32601, Message: "method not found: " + req.Method},
		}
	}
}

func (m *MockServer) handleInitialize(ctx context.Context, req RPCRequest) RPCResponse {
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx, req)
	}
	return m.defaultInitialize(ctx, req)
}

func (m *MockServer) defaultInitialize(ctx context.Context, req RPCRequest) RPCResponse {
	_ = ctx
	var params struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if b, err := json.Marshal(req.Params); err == nil {
		_ = json.Unmarshal(b, &params)
	}

	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "mock-mcp-server",
			"version": "1.0.0",
		},
	}

	b, _ := json.Marshal(result)
	m.initialized.Store(true)

	return RPCResponse{JSONRPC: "2.0", ID: mustMarshalID(req.ID), Result: b}
}

func (m *MockServer) handleListTools(ctx context.Context, req RPCRequest) RPCResponse {
	tools, err := m.ListToolsFunc(ctx)
	if err != nil {
		return RPCResponse{
			JSONRPC: "2.0",
			ID:      mustMarshalID(req.ID),
			Error:   &RPCError{Code: -32000, Message: err.Error()},
		}
	}

	result := map[string]interface{}{"tools": tools}
	b, _ := json.Marshal(result)

	return RPCResponse{JSONRPC: "2.0", ID: mustMarshalID(req.ID), Result: b}
}

func (m *MockServer) defaultListTools(ctx context.Context) ([]Tool, error) {
	_ = ctx
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Tool, len(m.tools))
	copy(result, m.tools)
	return result, nil
}

func (m *MockServer) handleCallTool(ctx context.Context, req RPCRequest) RPCResponse {
	if !m.initialized.Load() {
		return RPCResponse{
			JSONRPC: "2.0",
			ID:      mustMarshalID(req.ID),
			Error:   &RPCError{Code: -32000, Message: "not initialized"},
		}
	}

	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if b, err := json.Marshal(req.Params); err == nil {
		_ = json.Unmarshal(b, &params)
	}

	if params.Name == "" {
		return RPCResponse{
			JSONRPC: "2.0",
			ID:      mustMarshalID(req.ID),
			Error:   &RPCError{Code: -32602, Message: "missing tool name"},
		}
	}

	result, err := m.CallToolFunc(ctx, params.Name, params.Arguments)
	if err != nil {
		return RPCResponse{
			JSONRPC: "2.0",
			ID:      mustMarshalID(req.ID),
			Error:   &RPCError{Code: -32000, Message: err.Error()},
		}
	}

	b, _ := json.Marshal(result)

	m.mu.Lock()
	m.callCount[params.Name]++
	m.lastCallArgs[params.Name] = params.Arguments
	m.mu.Unlock()

	return RPCResponse{JSONRPC: "2.0", ID: mustMarshalID(req.ID), Result: b}
}

func (m *MockServer) defaultCallTool(ctx context.Context, name string, args map[string]interface{}) (ToolResult, error) {
	_ = ctx

	switch name {
	case "echo":
		msg, _ := args["message"].(string)
		return ToolResult{
			Content: []ToolContent{{Type: "text", Text: msg}},
		}, nil
	case "add":
		a, _ := args["a"].(float64)
		b, _ := args["b"].(float64)
		return ToolResult{
			Content: []ToolContent{{Type: "text", Text: fmt.Sprintf("%v", a+b)}},
		}, nil
	default:
		return ToolResult{}, fmt.Errorf("unknown tool: %s", name)
	}
}

func (m *MockServer) handlePing(ctx context.Context, req RPCRequest) RPCResponse {
	_ = ctx
	if !m.initialized.Load() {
		return RPCResponse{
			JSONRPC: "2.0",
			ID:      mustMarshalID(req.ID),
			Error:   &RPCError{Code: -32000, Message: "not initialized"},
		}
	}
	return RPCResponse{JSONRPC: "2.0", ID: mustMarshalID(req.ID), Result: json.RawMessage(`{}`)}
}

func (m *MockServer) AddTool(tool Tool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tools = append(m.tools, tool)
}

func (m *MockServer) GetCallCount(tool string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[tool]
}

func (m *MockServer) GetLastCallArgs(tool string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.lastCallArgs[tool]
	if args == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range args {
		result[k] = v
	}
	return result
}

func (m *MockServer) IsInitialized() bool {
	return m.initialized.Load()
}

func (m *MockServer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.initialized.Store(false)
	m.callCount = make(map[string]int)
	m.lastCallArgs = make(map[string]map[string]interface{})
}

type MockTransport struct {
	server    *MockServer
	closed    atomic.Bool
	callDelay int
}

func NewMockTransport(server *MockServer) *MockTransport {
	return &MockTransport{server: server}
}

func (t *MockTransport) SetCallDelay(ms int) {
	t.callDelay = ms
}

func (t *MockTransport) Call(ctx context.Context, req RPCRequest) (RPCResponse, error) {
	if t.closed.Load() {
		return RPCResponse{}, errors.New("transport closed")
	}

	if t.callDelay > 0 {
		select {
		case <-ctx.Done():
			return RPCResponse{}, ctx.Err()
		case <-time.After(time.Duration(t.callDelay) * time.Millisecond):
		}
	}

	return t.server.HandleRequest(ctx, req), nil
}

func (t *MockTransport) Notify(ctx context.Context, notif RPCNotification) error {
	if t.closed.Load() {
		return errors.New("transport closed")
	}
	req := RPCRequest{
		JSONRPC: notif.JSONRPC,
		Method:  notif.Method,
		Params:  notif.Params,
	}
	t.server.HandleRequest(ctx, req)
	return nil
}

func (t *MockTransport) Close() error {
	t.closed.Store(true)
	return nil
}

func (t *MockTransport) IsClosed() bool {
	return t.closed.Load()
}
