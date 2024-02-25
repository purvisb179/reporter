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

// ValidateAudience checks if the aud claim contains the clientID or if azp matches the clientID
func (s *OIDCService) ValidateAudience(claims struct {
	Audience interface{} `json:"aud"`
	Azp      string      `json:"azp"`
}) bool {
	switch aud := claims.Audience.(type) {
	case string:
		return aud == s.ClientId || claims.Azp == s.ClientId
	case []interface{}:
		for _, a := range aud {
			if str, ok := a.(string); ok && (str == s.ClientId || claims.Azp == s.ClientId) {
				return true
			}
		}
	}
	return false
}
