package oauth2

import "context"

type Provider interface {
	GetAuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*Token, error)
	VerifyIDToken(ctx context.Context, token *Token) (*IDToken, error)
}

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       int64
}

type IDToken struct {
	claims map[string]interface{}
}

func (t *IDToken) Claims(v interface{}) error {
	// Implement claims unmarshaling
	return nil
}
