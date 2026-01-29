package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"nhooyr.io/websocket"
	"pryx-core/internal/bus"
)

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow all origins for local dev
	})
	if err != nil {
		// Log error to stdout for now
		return
	}
	defer c.Close(websocket.StatusInternalError, "internal error")

	query := r.URL.Query()
	surface := strings.TrimSpace(query.Get("surface"))
	sessionFilter := strings.TrimSpace(query.Get("session_id"))
	eventFilters := query["event"]

	var topics []bus.EventType
	for _, ev := range eventFilters {
		ev = strings.TrimSpace(ev)
		if ev == "" {
			continue
		}
		topics = append(topics, bus.EventType(ev))
	}

	var events <-chan bus.Event
	var cancel func()
	if len(topics) == 0 {
		events, cancel = s.bus.Subscribe()
	} else {
		events, cancel = s.bus.Subscribe(topics...)
	}
	defer cancel()

	ctx := r.Context()

	s.bus.Publish(bus.NewEvent(bus.EventTraceEvent, sessionFilter, map[string]interface{}{
		"kind":        "ws.connected",
		"remote_addr": r.RemoteAddr,
		"surface":     surface,
	}))

	// Writer goroutine: Listen to bus, write to WS
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case evt, ok := <-events:
				if !ok {
					return
				}
				if sessionFilter != "" && evt.SessionID != sessionFilter {
					continue
				}
				bytes, err := json.Marshal(evt)
				if err != nil {
					continue
				}
				if err := c.Write(ctx, websocket.MessageText, bytes); err != nil {
					return
				}
			}
		}
	}()

	// Reader loop: Keep connection alive and handle incoming messages
	for {
		msgType, data, err := c.Read(ctx)
		if err != nil {
			break
		}
		if msgType != websocket.MessageText {
			continue
		}

		// Parse generic message structure
		in := struct {
			Event      string                 `json:"event"`
			Type       string                 `json:"type"`
			SessionID  string                 `json:"session_id"`
			Payload    map[string]interface{} `json:"payload"`
			ApprovalID string                 `json:"approval_id"`
			Approved   bool                   `json:"approved"`
		}{}
		if err := json.Unmarshal(data, &in); err != nil {
			continue
		}

		// Handle different message types
		eventType := in.Event
		if eventType == "" {
			eventType = in.Type
		}

		switch eventType {
		case "approval.resolve":
			if strings.TrimSpace(in.ApprovalID) != "" {
				_ = s.mcp.ResolveApproval(in.ApprovalID, in.Approved)
			}
		case "chat.send":
			if in.Payload != nil && in.Payload["content"] != nil {
				// Publish chat request for Agent to handle
				s.bus.Publish(bus.NewEvent(bus.EventChatRequest, sessionFilter, in.Payload))
			}
		}
	}

	s.bus.Publish(bus.NewEvent(bus.EventTraceEvent, sessionFilter, map[string]interface{}{
		"kind":        "ws.disconnected",
		"remote_addr": r.RemoteAddr,
		"surface":     surface,
	}))

	c.Close(websocket.StatusNormalClosure, "")
}
