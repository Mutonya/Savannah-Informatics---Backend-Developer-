package oauth2

import "context"

type MockProvider struct {
	GetAuthCodeURLFunc func(state string) string
	ExchangeFunc       func(ctx context.Context, code string) (*Token, error)
	VerifyIDTokenFunc  func(ctx context.Context, token *Token) (*IDToken, error)
}

func (m *MockProvider) GetAuthCodeURL(state string) string {
	return m.GetAuthCodeURLFunc(state)
}

func (m *MockProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	return m.ExchangeFunc(ctx, code)
}

func (m *MockProvider) VerifyIDToken(ctx context.Context, token *Token) (*IDToken, error) {
	return m.VerifyIDTokenFunc(ctx, token)
}
