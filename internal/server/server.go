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
	if env.MobileBackendURL != "" && env.MobileBackendKey != "" {
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
		r.Get("/mobile-api/wifi/get-wifi-plans", mobileBackendProxy)
		r.Post("/mobile-api/delete-account/has-deleted-account", mobileBackendProxy)
	}

	// proxy for `/network-api` requests to network backend
	if env.NetworkBackendURL != "" && env.NetworkBackendKey != "" {
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
	if env.DashboardBackendURL != "" && env.DashboardBackendKey != "" {
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
		dashboardBackendProxy := proxy.NewProxy(proxy.ProxyOptions{
			Target:           env.DashboardBackendURL,
			StripPrefix:      "/dashboard",
			Headers:          map[string]string{"X-API-Key": env.DashboardBackendKey},
			DisableForwarded: true,
			OverrideHost:     hostFromURL,
		})

		r.Get("/dashboard/", dashboardBackendProxy)
		r.Post("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Put("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
		r.Delete("/dashboard/", dashboardBackendProxy, gwmiddleware.KeycloakAuth(keycloakConfig))
	}

	// Keycloak config for admin panel routes
	adminKeycloakConfig := gwmiddleware.KeycloakAuthConfig{
		KeycloakUrl:   env.KeycloakUrl,
		KeycloakRealm: env.KeycloakAdminRealm,
		ClientID:      env.KeycloakAdminClientID,
		ClientSecret:  env.KeycloakAdminClientSecret,
	}

	// Proxy for /admin requests from admin panel to dashboard backend
	dashboardBackendAdminProxy := proxy.NewProxy(proxy.ProxyOptions{
		Target:           env.DashboardBackendURL,
		StripPrefix:      "/admin",
		Headers:          map[string]string{"X-API-Key": env.DashboardBackendAdminKey, "Authorization": fmt.Sprintf("Bearer %s", env.DashboardBackendAdminKey)},
		DisableForwarded: true,
		OverrideHost:     "",
	})

	r.Get("/admin/api/nfnodes", dashboardBackendAdminProxy, gwmiddleware.KeycloakAuth(adminKeycloakConfig))
	r.Get("/admin/api/nfnodes/{id}", dashboardBackendAdminProxy, gwmiddleware.KeycloakAuth(adminKeycloakConfig))
	r.Get("/admin/api/rewards-per-epoches", dashboardBackendAdminProxy, gwmiddleware.KeycloakAuth(adminKeycloakConfig))
	r.Get("/admin/api/rewards-per-epoches/{id}", dashboardBackendAdminProxy, gwmiddleware.KeycloakAuth(adminKeycloakConfig))
	r.Get("/admin/api/transaction-trackers", dashboardBackendAdminProxy, gwmiddleware.KeycloakAuth(adminKeycloakConfig))
	r.Get("/admin/api/transaction-trackers/{id}", dashboardBackendAdminProxy, gwmiddleware.KeycloakAuth(adminKeycloakConfig))

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
