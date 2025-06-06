package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// TestProvider is a mock OAuth provider for testing
type TestProvider struct {
	ExpectedState string
	ValidCode     string
	Token         *Token
	UserID        string
	UserEmail     string
}

// NewTestProvider creates a new test OAuth provider
func NewTestProvider() *TestProvider {
	return &TestProvider{
		ValidCode: "valid_test_code",
		Token:     &Token{AccessToken: "test_access_token"},
		UserID:    "test_user_123",
		UserEmail: "test@example.com",
	}
}

func (p *TestProvider) GetAuthCodeURL(state string) string {
	p.ExpectedState = state
	return "https://test-oauth-provider.com/auth?state=" + state
}

func (p *TestProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	if code != p.ValidCode {
		return nil, errors.New("invalid code")
	}
	return p.Token, nil
}

func (p *TestProvider) VerifyIDToken(ctx context.Context, token *Token) (*IDToken, error) {
	if token.AccessToken != p.Token.AccessToken {
		return nil, errors.New("invalid token")
	}

	return &IDToken{
		claims: map[string]interface{}{
			"sub":   p.UserID,
			"email": p.UserEmail,
			"name":  "Test User",
		},
	}, nil
}

// generateRandomString creates a random string for testing
func (p *TestProvider) generateRandomString() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
