package infra

import (
	"errors"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type GatewayEnvironment struct {
	AppEnv               string
	Port                 int
	KeycloakUrl          string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
	IdpServiceURL        string
	IdpServiceKey        string
	MobileBackendURL     string
	MobileBackendKey     string
	NetworkBackendURL    string
	NetworkBackendKey    string
}

func LoadEnvironment() (GatewayEnvironment, error) {
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		return GatewayEnvironment{}, errors.New("APP_ENV not set")
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return GatewayEnvironment{}, errors.New("PORT not set")
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return GatewayEnvironment{}, errors.New("PORT must be an integer")
	}

	keycloakUrl := strings.TrimSpace(os.Getenv("KEYCLOAK_URL"))
	if keycloakUrl == "" {
		return GatewayEnvironment{}, errors.New("KEYCLOAK_URL not set")
	}

	keycloakRealm := strings.TrimSpace(os.Getenv("KEYCLOAK_REALM"))
	if keycloakRealm == "" {
		return GatewayEnvironment{}, errors.New("KEYCLOAK_REALM not set")
	}

	keycloakClientID := strings.TrimSpace(os.Getenv("KEYCLOAK_CLIENT_ID"))
	if keycloakClientID == "" {
		return GatewayEnvironment{}, errors.New("KEYCLOAK_CLIENT_ID not set")
	}

	keycloakClientSecret := strings.TrimSpace(os.Getenv("KEYCLOAK_CLIENT_SECRET"))
	if keycloakClientSecret == "" {
		return GatewayEnvironment{}, errors.New("KEYCLOAK_CLIENT_SECRET not set")
	}

	idpServiceURL := strings.TrimSpace(os.Getenv("IDP_SERVICE_URL"))
	if idpServiceURL == "" {
		return GatewayEnvironment{}, errors.New("IDP_SERVICE_URL not set")
	}

	idpServiceKey := strings.TrimSpace(os.Getenv("IDP_SERVICE_KEY"))
	if idpServiceKey == "" {
		return GatewayEnvironment{}, errors.New("IDP_SERVICE_KEY not set")
	}

	mobileBackendURL := strings.TrimSpace(os.Getenv("MOBILE_BACKEND_URL"))
	mobileBackendKey := strings.TrimSpace(os.Getenv("MOBILE_BACKEND_KEY"))

	networkBackendURL := strings.TrimSpace(os.Getenv("NETWORK_BACKEND_URL"))
	networkBackendKey := strings.TrimSpace(os.Getenv("NETWORK_BACKEND_KEY"))

	return GatewayEnvironment{
		AppEnv:               appEnv,
		Port:                 portInt,
		KeycloakUrl:          keycloakUrl,
		KeycloakRealm:        keycloakRealm,
		KeycloakClientID:     keycloakClientID,
		KeycloakClientSecret: keycloakClientSecret,
		IdpServiceURL:        idpServiceURL,
		IdpServiceKey:        idpServiceKey,
		MobileBackendURL:     mobileBackendURL,
		MobileBackendKey:     mobileBackendKey,
		NetworkBackendURL:    networkBackendURL,
		NetworkBackendKey:    networkBackendKey,
	}, nil
}
