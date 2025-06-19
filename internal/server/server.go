package server

import (
	"fmt"
	"net/http"
	"net/url"
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

	// Proxy `/mobile-api` requests to mobile backend
	logger.Info("About to register proxy for mobile-api requests")
	if env.MobileBackendURL != "" && env.MobileBackendKey != "" {
		msg := fmt.Sprintf("Env url %s", env.MobileBackendURL)
		logger.Info(msg)

		// Extract host from MobileBackendURL for OverrideHost
		parsedURL, err := url.Parse(env.MobileBackendURL)
		if err != nil {
			logger.Error("Failed to parse MobileBackendURL", zap.Error(err))
			return nil, err
		}
		hostFromURL := parsedURL.Host

		mobileBackendProxy := proxy.NewProxy(proxy.ProxyOptions{
			Target:           env.MobileBackendURL,
			StripPrefix:      "/mobile-api",
			Headers:          map[string]string{"X-API-Key": env.MobileBackendKey},
			DisableForwarded: true,
			OverrideHost:     hostFromURL,
		})

		r.Handle("/mobile-api/", mobileBackendProxy)

		// r.Get("/mobile-api/", mobileBackendProxy)
		// r.Post("/mobile-api/", mobileBackendProxy)
		// r.Put("/mobile-api/", mobileBackendProxy)
		// r.Delete("/mobile-api/", mobileBackendProxy)
	}

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
