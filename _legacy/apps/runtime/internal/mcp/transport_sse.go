package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SSETransport struct {
	baseURL      string
	headers      map[string]string
	client       *http.Client
	sseEndpoint  string
	postEndpoint string

	mu          sync.RWMutex
	conn        *sseConnection
	pending     map[string]chan RPCResponse
	closed      bool
	lastEventID string
}

type sseConnection struct {
	resp   *http.Response
	reader *bufio.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSSETransport(endpointURL string, headers map[string]string) *SSETransport {
	return &SSETransport{
		baseURL:      strings.TrimRight(endpointURL, "/"),
		headers:      headers,
		client:       &http.Client{Timeout: 30 * time.Second},
		pending:      make(map[string]chan RPCResponse),
		sseEndpoint:  "/sse",
		postEndpoint: "/message",
	}
}

func (t *SSETransport) SetEndpoints(ssePath, postPath string) {
	t.sseEndpoint = ssePath
	t.postEndpoint = postPath
}

func (t *SSETransport) Connect(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		return nil
	}

	url := t.baseURL + t.sseEndpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create SSE request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	for k, v := range t.headers {
		if strings.TrimSpace(k) != "" {
			req.Header.Set(k, v)
		}
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("connect SSE: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("SSE connection failed: %d - %s", resp.StatusCode, string(body))
	}

	connCtx, cancel := context.WithCancel(context.Background())
	t.conn = &sseConnection{
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
		ctx:    connCtx,
		cancel: cancel,
	}

	go t.readLoop()

	return nil
}

func (t *SSETransport) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true

	for _, ch := range t.pending {
		close(ch)
	}
	t.pending = make(map[string]chan RPCResponse)

	if t.conn != nil {
		t.conn.cancel()
		t.conn.resp.Body.Close()
		t.conn = nil
	}
	t.mu.Unlock()

	return nil
}

func (t *SSETransport) Call(ctx context.Context, req RPCRequest) (RPCResponse, error) {
	if err := t.Connect(ctx); err != nil {
		return RPCResponse{}, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	var idRaw json.RawMessage
	if req.ID != nil {
		idRaw, _ = json.Marshal(req.ID)
	}
	key := idKey(idRaw)
	if key == "" {
		return RPCResponse{}, errors.New("request ID is required")
	}

	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return RPCResponse{}, errors.New("transport closed")
	}
	ch := make(chan RPCResponse, 1)
	t.pending[key] = ch
	t.mu.Unlock()

	url := t.baseURL + t.postEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.mu.Lock()
		delete(t.pending, key)
		t.mu.Unlock()
		return RPCResponse{}, fmt.Errorf("create POST request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range t.headers {
		if strings.TrimSpace(k) != "" {
			httpReq.Header.Set(k, v)
		}
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		t.mu.Lock()
		delete(t.pending, key)
		t.mu.Unlock()
		return RPCResponse{}, fmt.Errorf("send request: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.mu.Lock()
		delete(t.pending, key)
		t.mu.Unlock()
		return RPCResponse{}, fmt.Errorf("request failed: %d", resp.StatusCode)
	}

	select {
	case <-ctx.Done():
		t.mu.Lock()
		delete(t.pending, key)
		t.mu.Unlock()
		return RPCResponse{}, ctx.Err()
	case rpcResp, ok := <-ch:
		if !ok {
			return RPCResponse{}, errors.New("transport closed")
		}
		return rpcResp, nil
	}
}

func (t *SSETransport) Notify(ctx context.Context, notif RPCNotification) error {
	if err := t.Connect(ctx); err != nil {
		return err
	}

	body, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}

	url := t.baseURL + t.postEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create POST request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range t.headers {
		if strings.TrimSpace(k) != "" {
			httpReq.Header.Set(k, v)
		}
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	resp.Body.Close()

	return nil
}

func (t *SSETransport) readLoop() {
	var currentData []string

	for {
		t.mu.RLock()
		if t.closed || t.conn == nil {
			t.mu.RUnlock()
			return
		}
		reader := t.conn.reader
		t.mu.RUnlock()

		line, err := reader.ReadString('\n')
		if err != nil {
			t.handleError(err)
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			t.dispatchEvent(currentData)
			currentData = nil
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			currentData = append(currentData, data)
		} else if strings.HasPrefix(line, "id:") {
			t.mu.Lock()
			t.lastEventID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
			t.mu.Unlock()
		}
	}
}

func (t *SSETransport) dispatchEvent(dataLines []string) {
	if len(dataLines) == 0 {
		return
	}

	payload := strings.Join(dataLines, "\n")
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return
	}

	var resp RPCResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		return
	}

	key := idKey(resp.ID)
	if key == "" {
		return
	}

	t.mu.Lock()
	ch, ok := t.pending[key]
	if ok {
		delete(t.pending, key)
	}
	t.mu.Unlock()

	if ok {
		ch <- resp
		close(ch)
	}
}

func (t *SSETransport) handleError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return
	}

	for _, ch := range t.pending {
		close(ch)
	}
	t.pending = make(map[string]chan RPCResponse)

	if t.conn != nil {
		t.conn.resp.Body.Close()
		t.conn = nil
	}

	_ = err
}

func (t *SSETransport) IsConnected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.conn != nil && !t.closed
}
