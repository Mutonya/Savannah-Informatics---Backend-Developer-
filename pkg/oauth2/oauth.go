package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OAuthProvider interface {
	GetAuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}

type OIDCProvider struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
	ctx      context.Context
}

func NewOIDCProvider(ctx context.Context, clientID, clientSecret, redirectURL, providerURL string) (*OIDCProvider, error) {
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	return &OIDCProvider{
		config:   config,
		verifier: verifier,
		ctx:      ctx,
	}, nil
}

func (p *OIDCProvider) GetAuthCodeURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *OIDCProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

func (p *OIDCProvider) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in token response")
	}

	return p.verifier.Verify(ctx, rawIDToken)
}

func (p *OIDCProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (map[string]interface{}, error) {
	userInfoURL := strings.TrimSuffix(p.config.Endpoint.AuthURL, "/auth") + "/userinfo"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal userinfo: %w", err)
	}

	return userInfo, nil
}
