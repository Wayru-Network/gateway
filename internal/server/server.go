package server

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Wayru-Network/gateway/internal/infra"
	gwmiddleware "github.com/Wayru-Network/gateway/pkg/middleware"
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
	keycloakConfig := gwmiddleware.KeycloakAuthConfig{
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
	r.Get("/idp/profiles/token", idpProxy, gwmiddleware.KeycloakAuth(keycloakConfig))

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

		// Proxy for Socket.IO connections
		socketIOProxy := proxy.NewProxy(proxy.ProxyOptions{
			Target:           env.MobileBackendURL,
			StripPrefix:      "/ws-mobile-api", // No strip prefix for socket.io
			Headers:          map[string]string{"X-API-Key": env.MobileBackendKey},
			DisableForwarded: false,
			OverrideHost:     "",
		})
		r.Handle("/ws-mobile-api/socket.io/", socketIOProxy)

		// Keycloak config for any route that needs keycloak middleware auth
		keycloakConfig := gwmiddleware.KeycloakAuthConfig{
			KeycloakUrl:   env.KeycloakUrl,
			KeycloakRealm: env.KeycloakRealm,
			ClientID:      env.KeycloakClientID,
			ClientSecret:  env.KeycloakClientSecret,
		}

		mobileBackendProxy := proxy.NewProxy(proxy.ProxyOptions{
			Target:           env.MobileBackendURL,
			StripPrefix:      "/mobile-api",
			Headers:          map[string]string{"X-API-Key": env.MobileBackendKey},
			DisableForwarded: true,
			OverrideHost:     hostFromURL,
		})

		r.Get("/mobile-api/", mobileBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Post("/mobile-api/", mobileBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Put("/mobile-api/", mobileBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Delete("/mobile-api/", mobileBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))

		// define public endpoints
		r.Get("/mobile-api/esim/bundles", mobileBackendProxy)
	}

	// proxy for `/network-api` requests to network backend
	logger.Info("About to register proxy for network-api requests")
	if env.NetworkBackendURL != "" && env.NetworkBackendKey != "" {
		msg := fmt.Sprintf("Env url %s", env.NetworkBackendURL)
		logger.Info(msg)

		// Extract host from NetworkBackendURL for OverrideHost
		parsedURL, err := url.Parse(env.NetworkBackendURL)
		if err != nil {
			logger.Error("Failed to parse NetworkBackendURL", zap.Error(err))
			return nil, err
		}
		hostFromURL := parsedURL.Host

		// keycloak config for any route that needs keycloak middleware auth
		keycloakConfig := gwmiddleware.KeycloakAuthConfig{
			KeycloakUrl:   env.KeycloakUrl,
			KeycloakRealm: env.KeycloakRealm,
			ClientID:      env.KeycloakClientID,
			ClientSecret:  env.KeycloakClientSecret,
		}
		networkBackendProxy := proxy.NewProxy(proxy.ProxyOptions{
			Target:           env.NetworkBackendURL,
			StripPrefix:      "/network-api",
			Headers:          map[string]string{"X-API-Key": env.NetworkBackendKey},
			DisableForwarded: true,
			OverrideHost:     hostFromURL,
		})

		r.Get("/network-api/", networkBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Post("/network-api/", networkBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Put("/network-api/", networkBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Delete("/network-api/", networkBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))

	}

	// proxy for `/dashboard-api` requests to dashboard backend
	logger.Info("About to register proxy for dashboard-api requests")
	if env.DashboardBackendURL != "" && env.DashboardBackendKey != "" {
		msg := fmt.Sprintf("Env url %s", env.DashboardBackendURL)
		logger.Info(msg)

		// Extract host from NetworkBackendURL for OverrideHost
		parsedURL, err := url.Parse(env.NetworkBackendURL)
		if err != nil {
			logger.Error("Failed to parse NetworkBackendURL", zap.Error(err))
			return nil, err
		}
		hostFromURL := parsedURL.Host

		// keycloak config for any route that needs keycloak middleware auth
		keycloakConfig := gwmiddleware.KeycloakAuthConfig{
			KeycloakUrl:   env.KeycloakUrl,
			KeycloakRealm: env.KeycloakRealm,
			ClientID:      env.KeycloakClientID,
			ClientSecret:  env.KeycloakClientSecret,
		}
		logger.Info("dashboard backend url", zap.String("url", env.DashboardBackendURL))
		dashboardBackendProxy := proxy.NewProxy(proxy.ProxyOptions{
			Target:           env.DashboardBackendURL,
			StripPrefix:      "/dashboard",
			Headers:          map[string]string{"X-API-Key": env.DashboardBackendKey},
			DisableForwarded: true,
			OverrideHost:     hostFromURL,
		})

		r.Get("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Post("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Put("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Delete("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
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
