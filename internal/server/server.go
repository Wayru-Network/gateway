package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Wayru-Network/gateway/internal/infra"
	"github.com/Wayru-Network/serve/middleware"
	"github.com/Wayru-Network/serve/proxy"
	"github.com/Wayru-Network/serve/router"
	"go.uber.org/zap"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewServer(env infra.GatewayEnvironment) (*http.Server, error) {
	logger := zap.L()
	infra.ConfigureServeLogger(logger)
	r := router.NewRouter(middleware.RequestLogger())

	// Keycloak config for any route that needs keycloak middleware auth
	keycloakConfig := middleware.KeycloakAuthConfig{
		KeycloakUrl:   env.KeycloakUrl,
		KeycloakRealm: env.KeycloakRealm,
		ClientID:      env.KeycloakClientID,
		ClientSecret:  env.KeycloakClientSecret,
	}

	// IDP proxy instance
	idpProxy := proxy.NewProxy(proxy.ProxyOptions{
		Target:      env.IdpServiceURL,
		StripPrefix: "/idp",
		Headers:     map[string]string{"X-API-Key": env.IdpServiceKey},
	})

	// Proxy all `/idp` GET requests to idp api
	r.Get("/idp/", idpProxy)

	// Proxy for temporary idp token (more specific path match takes precedence)
	r.Get("/idp/profiles/token", idpProxy, middleware.KeycloakAuth(keycloakConfig))

	// Health
	r.Get("/health", health)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}, nil
}
