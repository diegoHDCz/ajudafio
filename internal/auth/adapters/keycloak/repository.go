package keycloak

import (
	"context"
	"log"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = "app-ajudafio"
	clientSecret = "HJt6rlZCK0W9smwde8NMW9Z37p5E7Hd4"
)

type KeycloakRepository struct {
	client  string
	baseURL string
}

func NewKeycloakRepository(baseURL string) *KeycloakRepository {
	return &KeycloakRepository{
		client:  clientID,
		baseURL: baseURL,
	}
}

func (k *KeycloakRepository) GetKeycloakConfig() (oauth2.Config, error) {
	ctx := context.Background()
	provicer, err := oidc.NewProvider(ctx, "http://localhost:8180/realms/ajudafio")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)

	}
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provicer.Endpoint(),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
	}

	return config, nil
}
