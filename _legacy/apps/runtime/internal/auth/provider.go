package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"pryx-core/internal/keychain"
)

// ProviderOAuth handles OAuth 2.0 flows for AI providers
type ProviderOAuth struct {
	keychain *keychain.Keychain
}

// NewProviderOAuth creates a new provider OAuth handler
func NewProviderOAuth(kc *keychain.Keychain) *ProviderOAuth {
	return &ProviderOAuth{keychain: kc}
}

// ProviderConfig holds OAuth configuration for a provider
type ProviderConfig struct {
	Name        string
	ClientID    string
	AuthURL     string
	TokenURL    string
	Scopes      []string
	PKCEEnabled bool
}

// ProviderConfigs defines OAuth configurations for supported providers
var ProviderConfigs = map[string]ProviderConfig{
	"google": {
		Name:        "Google",
		ClientID:    "93780524682-mq1q8n2e4k2q5d6p4v3n1s9c5o1p2q.apps.googleusercontent.com", // Replace with actual
		AuthURL:     "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:    "https://oauth2.googleapis.com/token",
		Scopes:      []string{"https://www.googleapis.com/auth/generative-language.retroactive"},
		PKCEEnabled: true,
	},
}

// StartOAuthFlow initiates OAuth flow with local callback server
func (p *ProviderOAuth) StartOAuthFlow(ctx context.Context, providerID string) (*TokenResponse, error) {
	config, ok := ProviderConfigs[providerID]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", providerID)
	}

	// Start local callback server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}
	defer listener.Close()

	callbackURL := fmt.Sprintf("http://localhost:%d/oauth/callback", listener.Addr().(*net.TCPAddr).Port)

	// Generate PKCE parameters
	var codeChallenge, codeVerifier string
	if config.PKCEEnabled {
		verifier, err := generateCodeVerifier()
		if err != nil {
			return nil, fmt.Errorf("failed to generate PKCE: %w", err)
		}
		codeVerifier = verifier
		codeChallenge = generateCodeChallenge(verifier)
	}

	state, err := generateRandomState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// Build authorization URL
	authURL, err := buildAuthURL(config, callbackURL, state, codeChallenge)
	if err != nil {
		return nil, fmt.Errorf("failed to build auth URL: %w", err)
	}

	// Start HTTP server to handle callback
	codeChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleOAuthCallback(w, r, state, codeChan, errorChan)
		}),
	}

	go server.Serve(listener)
	defer server.Close()

	// Open browser
	fmt.Printf("Opening browser for %s OAuth...\n", config.Name)
	fmt.Printf("URL: %s\n", authURL)

	// Wait for callback or timeout
	select {
	case code := <-codeChan:
		return p.exchangeCode(ctx, config, code, callbackURL, codeVerifier)
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(5 * time.Minute):
		return nil, errors.New("OAuth timeout - user did not complete authorization")
	}
}

// buildAuthURL constructs the OAuth authorization URL
func buildAuthURL(config ProviderConfig, redirectURI, state, codeChallenge string) (string, error) {
	params := url.Values{
		"client_id":     {config.ClientID},
		"redirect_uri":  {redirectURI},
		"response_type": {"code"},
		"scope":         {joinScopes(config.Scopes)},
		"state":         {state},
		"access_type":   {"offline"},
		"prompt":        {"consent"},
	}

	if config.PKCEEnabled && codeChallenge != "" {
		params.Set("code_challenge", codeChallenge)
		params.Set("code_challenge_method", "S256")
	}

	return config.AuthURL + "?" + params.Encode(), nil
}

// handleOAuthCallback processes the OAuth callback
func handleOAuthCallback(w http.ResponseWriter, r *http.Request, expectedState string, codeChan chan<- string, errorChan chan<- error) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		errorChan <- fmt.Errorf("OAuth error: %s", errorParam)
		http.Error(w, "Authorization failed: "+errorParam, http.StatusBadRequest)
		return
	}

	if state != expectedState {
		errorChan <- errors.New("state mismatch - possible CSRF attack")
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	if code == "" {
		errorChan <- errors.New("no authorization code received")
		http.Error(w, "No code", http.StatusBadRequest)
		return
	}

	codeChan <- code

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>OAuth Complete</title></head>
<body>
<h1>Authorization Successful!</h1>
<p>You can close this window and return to Pryx.</p>
<script>setTimeout(() => window.close(), 3000);</script>
</body>
</html>`))
}

// exchangeCode exchanges authorization code for tokens
func (p *ProviderOAuth) exchangeCode(ctx context.Context, config ProviderConfig, code, redirectURI, codeVerifier string) (*TokenResponse, error) {
	params := url.Values{
		"client_id":    {config.ClientID},
		"code":         {code},
		"grant_type":   {"authorization_code"},
		"redirect_uri": {redirectURI},
	}

	if config.PKCEEnabled && codeVerifier != "" {
		params.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.TokenURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// SaveTokens stores OAuth tokens in keychain
func (p *ProviderOAuth) SaveTokens(providerID string, tokens *TokenResponse) error {
	if p.keychain == nil {
		return errors.New("keychain not available")
	}

	prefix := "oauth_" + providerID + "_"

	if err := p.keychain.Set(prefix+"access", tokens.AccessToken); err != nil {
		return fmt.Errorf("failed to save access token: %w", err)
	}

	if tokens.RefreshToken != "" {
		if err := p.keychain.Set(prefix+"refresh", tokens.RefreshToken); err != nil {
			return fmt.Errorf("failed to save refresh token: %w", err)
		}
	}

	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	if err := p.keychain.Set(prefix+"expires", expiresAt.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("failed to save expiry: %w", err)
	}

	return nil
}

// GetToken retrieves access token from keychain
func (p *ProviderOAuth) GetToken(providerID string) (string, error) {
	if p.keychain == nil {
		return "", errors.New("keychain not available")
	}

	return p.keychain.Get("oauth_" + providerID + "_access")
}

// RefreshToken refreshes an expired access token
func (p *ProviderOAuth) RefreshToken(ctx context.Context, providerID string) error {
	config, ok := ProviderConfigs[providerID]
	if !ok {
		return fmt.Errorf("unsupported provider: %s", providerID)
	}

	refreshToken, err := p.keychain.Get("oauth_" + providerID + "_refresh")
	if err != nil {
		return fmt.Errorf("no refresh token available: %w", err)
	}

	params := url.Values{
		"client_id":     {config.ClientID},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.TokenURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refresh failed: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode refresh response: %w", err)
	}

	// Update stored tokens
	return p.SaveTokens(providerID, &tokenResp)
}

// IsTokenExpired checks if token needs refresh
func (p *ProviderOAuth) IsTokenExpired(providerID string) (bool, error) {
	expiresStr, err := p.keychain.Get("oauth_" + providerID + "_expires")
	if err != nil {
		return true, nil // Assume expired if not found
	}

	expiresAt, err := time.Parse(time.RFC3339, expiresStr)
	if err != nil {
		return true, nil
	}

	// Refresh 5 minutes before expiry
	return time.Until(expiresAt) < 5*time.Minute, nil
}

// generateCodeVerifier generates a PKCE code verifier
func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge creates S256 code challenge from verifier
func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// joinScopes joins OAuth scopes with spaces
func joinScopes(scopes []string) string {
	result := ""
	for i, s := range scopes {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}
