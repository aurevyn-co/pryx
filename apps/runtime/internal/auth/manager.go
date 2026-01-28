package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"pryx-core/internal/bus"
	"pryx-core/internal/keychain"
	"pryx-core/internal/store"
)

const (
	OAuthDeviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"
	OAuthAuthorizationCodeGrantType = "authorization_code"
	TokenValiditySeconds            = 3600 // 1 hour
	TokenRefreshBufferSeconds     = 300   // 5 minutes before refresh
)

type OAuthProvider struct {
	Name           string `json:"name"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	AuthURL       string `json:"auth_url"`
	TokenURL      string `json:"token_url"`
	Scopes        []string `json:"scopes"`
}

type OAuthState struct {
	State        string    `json:"state"`
	ProviderID   string    `json:"provider_id"`
	ClientID     string    `json:"client_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	RedirectURI  string    `json:"redirect_uri,omitempty"`
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope       string `json:"scope"`
	ProviderID   string `json:"provider_id"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope       string `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type TokenConfig struct {
	ProviderID string `json:"provider_id"`
	Token      string `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Manager struct {
	config      *Config
	keychain    *keychain.Keychain
	store       *store.Store
	bus         *bus.Bus
}

type Config struct {
	OAuthProviders map[string]*OAuthProvider `json:"oauth_providers"`
}

func NewManager(config *Config, keychain *keychain.Keychain, store *store.Store, bus *bus.Bus) *Manager {
	return &Manager{
		config:   config,
		keychain: keychain,
		store:    store,
		bus:      bus,
	}
}

func (m *Manager) InitiateDeviceFlow(ctx context.Context, providerID string, redirectURI string) (*OAuthState, error) {
	provider, ok := m.config.OAuthProviders[providerID]
	if !ok {
		return nil, errors.New("provider not found")
	}

	state, err := m.createOAuthState(ctx, provider, redirectURI)
	if err != nil {
		return nil, err
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventOAuthFlowInitiated, providerID, map[string]interface{}{
			"provider_id": providerID,
			"state":      state.State,
			"redirect_uri": redirectURI,
		})
		m.bus.Publish(event)
	}

	// Save state to keychain
	err = m.keychain.Set(ctx, "oauth_state_"+state.State, state.State)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (m *Manager) createOAuthState(ctx context.Context, provider *OAuthProvider, redirectURI string) (*OAuthState, error) {
	stateBytes, err := generateRandomState()
	if err != nil {
		return nil, err
	}

	state := string(stateBytes)
	expiresAt := time.Now().Add(TokenValiditySeconds)

	return &OAuthState{
		State:        state,
		ProviderID:   provider.ClientID,
		ClientID:     provider.ClientID,
		ExpiresAt:    expiresAt,
		RedirectURI:  redirectURI,
	}, nil
}

func (m *Manager) PollDeviceAuth(ctx context.Context, state string) (*OAuthToken, error) {
	oauthState, err := m.getOAuthState(ctx, state)
	if err != nil {
		return nil, err
	}

	provider, ok := m.config.OAuthProviders[oauthState.ProviderID]
	if !ok {
		return nil, errors.New("provider not found")
	}

	// Poll for token using device code
	tokenURL := fmt.Sprintf("%s?client_id=%s&grant_type=%s&device_code=%s",
		provider.TokenURL,
		provider.ClientID,
		OAuthDeviceCodeGrantType,
		state.State,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn   int64  `json:"expires_in"`
		Scope       string `json:"scope"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	token := OAuthToken{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresAt:  expiresAt,
		Scope:      result.Scope,
		ProviderID: oauthState.ProviderID,
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventOAuthTokenReceived, state, map[string]interface{}{
			"provider_id": providerID,
			"token":       token.AccessToken,
			"expires_at": token.ExpiresAt,
		})
		m.bus.Publish(event)
	}

	return &token, nil
}

func (m *Manager) RefreshToken(ctx context.Context, providerID string) (*OAuthToken, error) {
	provider, ok := m.config.OAuthProviders[providerID]
	if !ok {
		return nil, errors.New("provider not found")
	}

	tokenConfig, err := m.keychain.Get(ctx, "oauth_token_"+providerID)
	if err != nil {
		return nil, err
	}

	var tokenConfig TokenConfig
	if err := json.Unmarshal([]byte(tokenConfig.Token), &tokenConfig); err != nil {
		return nil, err
	}

	if time.Now().Before(tokenConfig.ExpiresAt.Add(-TokenRefreshBufferSeconds)) {
		return m.refreshTokenFromKeychain(ctx, provider, tokenConfig.RefreshToken)
	}

	return &OAuthToken{
		AccessToken: tokenConfig.Token,
		TokenType:   tokenConfig.TokenType,
		ExpiresAt:  tokenConfig.ExpiresAt,
		ProviderID: providerID,
	}, nil
}

func (m *Manager) refreshTokenFromKeychain(ctx context.Context, provider *OAuthProvider, refreshToken string) (*OAuthToken, error) {
	if refreshToken == "" {
		return nil, errors.New("no refresh token available")
	}

	tokenURL := fmt.Sprintf("%s?client_id=%s&grant_type=refresh_token&refresh_token=%s",
		provider.TokenURL,
		provider.ClientID,
		refreshToken,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(tokenURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(TokenValiditySeconds)
	token := OAuthToken{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresAt:  expiresAt,
		Scope:      result.Scope,
		ProviderID: provider.Name,
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventOAuthTokenRefreshed, provider.Name, map[string]interface{}{
			"new_access_token": result.AccessToken,
			"expires_at":      token.ExpiresAt,
		})
		m.bus.Publish(event)
	}

	err = m.keychain.Set(ctx, "oauth_token_"+provider.Name, token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (m *Manager) SetManualToken(ctx context.Context, providerID string, token string) error {
	if token == "" {
		return errors.New("token cannot be empty")
	}

	expiresAt := time.Now().Add(TokenValiditySeconds)
	oauthToken := OAuthToken{
		AccessToken: token,
		TokenType:   "manual",
		ExpiresAt:  expiresAt,
		ProviderID: providerID,
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventOAuthTokenSet, providerID, map[string]interface{}{
			"token_type": "manual",
			"expires_at": oauthToken.ExpiresAt,
		})
		m.bus.Publish(event)
	}

	err = m.keychain.Set(ctx, "oauth_token_"+providerID, oauthToken)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetTokens(ctx context.Context) (map[string]*OAuthToken, error) {
	tokens := make(map[string]*OAuthToken)

	for providerID := range m.config.OAuthProviders {
		token, err := m.keychain.Get(ctx, "oauth_token_"+providerID)
		if err != nil {
			continue
		}

		var oauthToken OAuthToken
		if err := json.Unmarshal([]byte(token), &oauthToken); err != nil {
			continue
		}

		oauthToken.ProviderID = providerID
		tokens[providerID] = &oauthToken
	}

	return tokens, nil
}

func (m *Manager) ValidateToken(ctx context.Context, providerID string) error {
	token, err := m.keychain.Get(ctx, "oauth_token_"+providerID)
	if err != nil {
		return err
	}

	var oauthToken OAuthToken
	if err := json.Unmarshal([]byte(token), &oauthToken); err != nil {
		return err
	}

	if time.Now().After(oauthToken.ExpiresAt) {
		return errors.New("token has expired")
	}

	return nil
}

func (m *Manager) RevokeToken(ctx context.Context, providerID string) error {
	err := m.keychain.Delete(ctx, "oauth_token_"+providerID)
	if err != nil {
		return err
	}

	if m.bus != nil {
		event := bus.NewEvent(bus.EventOAuthTokenRevoked, providerID, nil)
		m.bus.Publish(event)
	}

	return nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) cleanupExpiredStates(ctx context.Context) error {
	states, err := m.keychain.ListByPrefix(ctx, "oauth_state_")
	if err != nil {
		return err
	}

	deletedCount := 0
	now := time.Now()

	for key, value := range states {
		if strings.HasPrefix(key, "oauth_state_") {
			var state OAuthState
			if err := json.Unmarshal([]byte(value), &state); err != nil {
				continue
			}

			if now.After(state.ExpiresAt) {
				err := m.keychain.Delete(ctx, key)
				if err == nil {
					deletedCount++
				}
			}
		}
	}

	return nil
}

func (m *Manager) GetOAuthConfig(ctx context.Context) (*OAuthConfig, error) {
	providers := make(map[string]*OAuthProvider)

	for providerID, provider := range m.config.OAuthProviders {
		providers[providerID] = &provider
	}

	return &OAuthConfig{OAuthProviders: providers}, nil
}
