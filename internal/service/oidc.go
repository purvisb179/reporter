package service

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"log"
)

// OIDCService encapsulates the OIDC provider and verifier.
type OIDCService struct {
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
	Config   oauth2.Config
	ctx      context.Context
	ClientId string
}

// NewOIDCService creates a new OIDCService with the specified parameters.
func NewOIDCService(ctx context.Context, providerURL, clientID, clientSecret string, redirectURL string) *OIDCService {
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatalf("Failed to get OIDC provider: %v", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &OIDCService{
		Provider: provider,
		Verifier: verifier,
		Config:   config,
		ctx:      ctx,
		ClientId: clientID,
	}
}

// VerifyToken verifies the OIDC token and returns the claims.
func (s *OIDCService) VerifyToken(token string) (*oidc.IDToken, error) {
	// Parse and verify the ID token payload
	idToken, err := s.Verifier.Verify(s.ctx, token)
	if err != nil {
		return nil, err
	}
	return idToken, nil
}
