package mesh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"pryx-core/internal/bus"
	"pryx-core/internal/config"
	"pryx-core/internal/keychain"
	"pryx-core/internal/store"
)

type Manager struct {
	cfg      *config.Config
	bus      *bus.Bus
	store    *store.Store
	keychain *keychain.Keychain
}

func NewManager(cfg *config.Config, b *bus.Bus, s *store.Store, kc *keychain.Keychain) *Manager {
	return &Manager{
		cfg:      cfg,
		bus:      b,
		store:    s,
		keychain: kc,
	}
}

func (m *Manager) Start(ctx context.Context) {
	// 1. Listen for local events to broadcast
	go m.listenForBroadcasts(ctx)

	// 2. Periodic sync
	go m.periodicSync(ctx)

	log.Println("Pryx Mesh Manager started")
}

func (m *Manager) listenForBroadcasts(ctx context.Context) {
	events, closer := m.bus.Subscribe(bus.EventSessionMessage, bus.EventSessionTyping)
	defer closer()

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			go m.broadcast(evt)
		}
	}
}

func (m *Manager) broadcast(evt bus.Event) {
	token, err := m.keychain.Get("cloud_access_token")
	if err != nil {
		return // Not logged in
	}

	payload, _ := json.Marshal(evt)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/sessions/broadcast", m.cfg.CloudAPIUrl), bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func (m *Manager) periodicSync(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.syncSessions(ctx)
		}
	}
}

func (m *Manager) syncSessions(ctx context.Context) {
	token, err := m.keychain.Get("cloud_access_token")
	if err != nil {
		return // Not logged in
	}

	// For now, we just fetch a list of sessions we might need to sync
	// In a real implementation, we'd compare versions or timestamps
	sessions, err := m.store.ListSessions()
	if err != nil {
		return
	}

	for _, s := range sessions {
		m.syncSession(ctx, s.ID, token)
	}
}

func (m *Manager) syncSession(ctx context.Context, sessionID string, token string) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/sessions/%s", m.cfg.CloudAPIUrl, sessionID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var cloudEvents []bus.Event
	if err := json.NewDecoder(resp.Body).Decode(&cloudEvents); err != nil {
		return
	}

	// Merge logic: publish events that we don't have locally
	// In this simplified version, we just re-publish them if they are newer
	// than our latest version (once we implement versioning properly)
	for _, _ = range cloudEvents {
		// Avoid infinite loops: only publish if it's from a different surface/device
		// and we are not the one who just broadcast it.
		// For now, we skip merge to avoid complexity in this step
		// But we could: m.bus.Publish(evt)
	}
}
