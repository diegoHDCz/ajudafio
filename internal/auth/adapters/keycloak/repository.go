package keycloak

import (
	"context"
	"log"
	"os"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type KeycloakRepository struct {
	client  string
	baseURL string
}

func NewKeycloakRepository(baseURL string) *KeycloakRepository {
	return &KeycloakRepository{
		client:  os.Getenv("CLIENT_ID"),
		baseURL: baseURL,
	}
}

func (k *KeycloakRepository) GetKeycloakConfig() (oauth2.Config, error) {
	ctx := context.Background()
	provicer, err := oidc.NewProvider(ctx, "http://localhost:8180/realms/ajudafio")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)

	}
	clientid := os.Getenv("CLIENT_ID")

	clientsecret := os.Getenv("CLIENT_SECRET")

	config := oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		Endpoint:     provicer.Endpoint(),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
	}

	return config, nil
}
