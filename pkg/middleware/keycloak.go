package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Wayru-Network/serve/router"
	"go.uber.org/zap"
)

type KeycloakAuthConfig struct {
	KeycloakUrl   string
	KeycloakRealm string
	ClientID      string
	ClientSecret  string
}

type KeycloakIntrospectResponse struct {
	Exp    int64  `json:"exp"` // Expires at
	Iat    int64  `json:"iat"` // Issued at
	Sub    string `json:"sub"` // User ID
	Active bool   `json:"active"`
}

func KeycloakAuth(config KeycloakAuthConfig) router.Middleware {
	logger := zap.L()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Keycloak access token from bearer
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				logger.Error("No authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check bearer format
			if !strings.HasPrefix(authorization, "Bearer ") {
				logger.Error("Invalid bearer format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Introspect
			bearer := authorization[len("Bearer "):]

			formData := url.Values{}
			formData.Set("client_id", config.ClientID)
			formData.Set("client_secret", config.ClientSecret)
			formData.Set("token", bearer)
			data := strings.NewReader(formData.Encode())

			url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", config.KeycloakUrl, config.KeycloakRealm)

			req, err := http.NewRequest("POST", url, data)
			if err != nil {
				logger.Error("Error creating POST request to introspect endpoint: " + err.Error())
				http.Error(w, "Bad Config", http.StatusInternalServerError)
				return
			}

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Error("Error sending POST request to introspect endpoint: " + err.Error())
				http.Error(w, "Bad Config or Comms", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			var instrospectResponse KeycloakIntrospectResponse
			err = json.NewDecoder(resp.Body).Decode(&instrospectResponse)
			if err != nil {
				logger.Error("Error processing introspect response: " + err.Error())
				http.Error(w, "Bad Processing", http.StatusInternalServerError)
				return
			}

			// Print introspect response
			logger.Info("Introspect response: token active = " + fmt.Sprintf("%t", instrospectResponse.Active))

			if !instrospectResponse.Active {
				logger.Error("Token is not active")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add sub as X-WAYRU-CONNECT-ID header
			r.Header.Set("X-WAYRU-CONNECT-ID", instrospectResponse.Sub)

			next.ServeHTTP(w, r)
		})
	}
}
