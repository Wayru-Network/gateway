# Gateway Service

A robust Go-based API Gateway microservice for routing and managing requests to multiple backend services. This gateway provides secure authentication, request proxying, and centralized routing for mobile, network, dashboard, and identity provider services.

## Features

🔐 **Keycloak Authentication**: JWT-based authentication using Keycloak token introspection

🚦 **Multi-Service Routing**: Intelligent routing to mobile, network, dashboard, and IDP backends

⚡ **WebSocket Support**: Full Socket.IO WebSocket proxy support for real-time connections

🛡️ **API Key Protection**: X-API-Key header injection for backend service authentication

🔄 **Request Logging**: Comprehensive request logging with structured logging (zap)

🐳 **Docker Support**: Containerized deployment with multi-stage Docker builds

☸️ **Kubernetes Ready**: Complete Kubernetes Helm charts for deployment

📊 **Health Monitoring**: Built-in health check endpoints for service monitoring

🔍 **Flexible Configuration**: Environment-based configuration for different deployment stages

## Tech Stack

- **Runtime**: Go 1.24.3+
- **Language**: Go
- **Authentication**: Keycloak (OpenID Connect)
- **Logging**: Zap (uber-go/zap)
- **Proxy**: Custom proxy implementation
- **Containerization**: Docker
- **Orchestration**: Kubernetes (Helm)

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.24.3 or higher
- Git configured with SSH access (for private repositories)
- Keycloak instance (for authentication)
- Docker (optional, for containerized deployment)
- Kubernetes cluster (optional, for Kubernetes deployment)

## Installation

### Clone the repository

```bash
git clone <repository-url>
cd gateway
```

### Set up Git for private repositories

This project may depend on private repositories. Configure Git and Go for private repository access:

1. Set up SSH access to GitHub (see [Git Setup Requirements](docs/git-setup.md))
2. Set the `GOPRIVATE` environment variable:
   ```bash
   export GOPRIVATE=github.com/Wayru-Network/*
   ```

### Install dependencies

```bash
go mod download
```

### Set up environment variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
APP_ENV=development
PORT=4050

# Keycloak Configuration
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=your_realm
KEYCLOAK_CLIENT_ID=your_client_id
KEYCLOAK_CLIENT_SECRET=your_client_secret

# IDP Service Configuration
IDP_SERVICE_URL=http://localhost:3000
IDP_SERVICE_KEY=your_idp_api_key

# Mobile Backend Configuration (optional)
MOBILE_BACKEND_URL=http://localhost:3001
MOBILE_BACKEND_KEY=your_mobile_backend_api_key

# Network Backend Configuration (optional)
NETWORK_BACKEND_URL=http://localhost:3002
NETWORK_BACKEND_KEY=your_network_backend_api_key

# Dashboard Backend Configuration (optional)
DASHBOARD_BACKEND_URL=http://localhost:3003
DASHBOARD_BACKEND_KEY=your_dashboard_backend_api_key
```

### Build the project

```bash
go build -o gateway ./cmd
```

Or using the justfile:

```bash
just build
```

### Start the server

```bash
./gateway
```

Or using the justfile:

```bash
just run
```

For development with auto-reload (requires `air`):

```bash
just watch
```

## Configuration

### Environment Variables

| Variable                 | Description                   | Required | Default |
| ------------------------ | ----------------------------- | -------- | ------- |
| `APP_ENV`                | Environment (local/dev/prod)  | Yes      | -       |
| `PORT`                   | Server port                   | Yes      | -       |
| `KEYCLOAK_URL`           | Keycloak server URL           | Yes      | -       |
| `KEYCLOAK_REALM`         | Keycloak realm name           | Yes      | -       |
| `KEYCLOAK_CLIENT_ID`     | Keycloak client ID            | Yes      | -       |
| `KEYCLOAK_CLIENT_SECRET` | Keycloak client secret        | Yes      | -       |
| `IDP_SERVICE_URL`        | Identity Provider service URL | Yes      | -       |
| `IDP_SERVICE_KEY`        | API key for IDP service       | Yes      | -       |
| `MOBILE_BACKEND_URL`     | Mobile backend service URL    | No       | -       |
| `MOBILE_BACKEND_KEY`     | API key for mobile backend    | No       | -       |
| `NETWORK_BACKEND_URL`    | Network backend service URL   | No       | -       |
| `NETWORK_BACKEND_KEY`    | API key for network backend   | No       | -       |
| `DASHBOARD_BACKEND_URL`  | Dashboard backend service URL | No       | -       |
| `DASHBOARD_BACKEND_KEY`  | API key for dashboard backend | No       | -       |

### Application Environments

The `APP_ENV` variable controls logging behavior:

- **local**: Development logging with human-readable output
- **dev**: Production logging with debug level enabled
- **prod**: Production logging with info level

## API Endpoints

### Health Check

**GET** `/health`

Returns `200 OK` if the service is healthy.

**Response:**

```
OK
```

### IDP Service Proxy

**GET** `/idp/*`

Proxies requests to the Identity Provider service. Requires API key authentication.

**GET** `/idp/profiles/token`

Proxies requests to IDP token endpoint. Requires Keycloak authentication.

### Mobile Backend Proxy

**GET/POST/PUT/DELETE** `/mobile-api/*`

Proxies requests to the mobile backend service. Most endpoints require Keycloak authentication.

**Public Endpoints** (no authentication required):

- `GET /mobile-api/esim/bundles`
- `GET /mobile-api/wifi/get-wifi-plans`
- `POST /mobile-api/delete-account/has-deleted-account`

**WebSocket Endpoint:**

- `GET /ws-mobile-api/socket.io/*` - Socket.IO WebSocket proxy

### Network Backend Proxy

**GET/POST/PUT/DELETE** `/network-api/*`

Proxies requests to the network backend service. Requires Keycloak authentication.

### Dashboard Backend Proxy

**GET** `/dashboard/*`

Proxies GET requests to the dashboard backend (no authentication).

**POST/PUT/DELETE** `/dashboard/*`

Proxies requests to the dashboard backend. Requires Keycloak authentication.

## Authentication

Most endpoints require Keycloak authentication via the `Authorization` header:

```
Authorization: Bearer <your_keycloak_token>
```

The gateway validates tokens by introspecting them with Keycloak. Upon successful validation, the user ID (`sub` claim) is added to the request as the `X-WAYRU-CONNECT-ID` header before forwarding to backend services.

## Docker Deployment

### Build Docker Image

```bash
docker build -t gateway .
```

**Note**: The Dockerfile expects SSH keys for private repository access. Ensure you have the necessary SSH setup or modify the Dockerfile for your use case.

### Run Container

```bash
docker run -p 4050:4050 --env-file .env gateway
```

### Docker Compose (Example)

```yaml
version: "3.8"
services:
  gateway:
    build: .
    ports:
      - "4050:4050"
    env_file:
      - .env
    depends_on:
      - keycloak
    restart: unless-stopped

  keycloak:
    image: quay.io/keycloak/keycloak:latest
    environment:
      KEYCLOAK_ADMIN: admin
      KEYCLOAK_ADMIN_PASSWORD: admin
    ports:
      - "8080:8080"
    command: start-dev
```

## Kubernetes Deployment

Kubernetes configuration files are available in the `deploy/chart/` directory. The project includes Helm charts for easy deployment.

### Prerequisites

- Kubernetes cluster
- Helm 3.x
- Keycloak instance accessible from the cluster

### Deploy

1. Update `deploy/chart/values.yaml` or environment-specific values files with your configuration
2. Install the chart:

```bash
helm install gateway ./deploy/chart -f ./deploy/chart/values-dev.yaml
```

### Update Configuration

```bash
helm upgrade gateway ./deploy/chart -f ./deploy/chart/values-prod.yaml
```

The Helm chart includes:

- Deployment with health probes
- Service configuration
- Ingress configuration
- CSI volume support for secrets (Azure Key Vault)

## Development

### Project Structure

```
gateway/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── infra/               # Infrastructure (env, logger)
│   └── server/              # Server setup and routing
├── pkg/
│   └── middleware/          # Custom middleware (Keycloak auth)
├── deploy/
│   └── chart/               # Kubernetes Helm charts
├── integration/             # Integration test scripts
├── docs/                    # Documentation
├── Dockerfile               # Docker configuration
├── go.mod                   # Go dependencies
├── justfile                 # Task runner commands
└── README.md                # This file
```

### Scripts

Using `just` (recommended):

- `just build` - Build the application
- `just run` - Run the application
- `just test` - Run tests
- `just watch` - Watch for changes and rebuild (requires `air`)

Or using Go directly:

- `go build -o gateway ./cmd` - Build the application
- `go run ./cmd/main.go` - Run the application
- `go test ./...` - Run tests

### Code Style

The project follows Go best practices:

- Use `gofmt` for code formatting
- Follow Go naming conventions
- Use structured logging with zap
- Handle errors explicitly
- Use context for request cancellation

### Testing

Run integration tests:

```bash
cd integration/keycloak_login && bash run.bash
cd integration/keycloak_introspect && bash run.bash
cd integration/get_profile && bash run.bash
```

## Security Considerations

- **Keycloak Secrets**: Never commit Keycloak client secrets to the repository. Use environment variables or secret management systems.
- **API Keys**: Store backend API keys securely using environment variables or Kubernetes secrets.
- **Private Repositories**: Configure SSH access properly for private Go module dependencies.
- **Token Validation**: All tokens are validated via Keycloak introspection before forwarding requests.
- **HTTPS**: Use HTTPS in production environments.
- **CORS**: Configure CORS appropriately if needed for your frontend applications.

## Error Handling

The gateway handles errors gracefully:

- **401 Unauthorized**: Invalid or missing authentication token
- **500 Internal Server Error**: Configuration errors or communication failures with Keycloak/backends

All errors are logged using structured logging for debugging and monitoring.

## Contributing

This project is now open source. Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

**Important**: This project is now open source and maintained by the community. WAYRU no longer exists and will not provide support for this repository. For issues, questions, or contributions, please use the GitHub Issues section.

💙 **Farewell Message**

With gratitude and love, we say goodbye.

WAYRU is closing its doors, but we are leaving these repositories open and free for the community.

May they continue to inspire builders, dreamers, and innovators.

With love, WAYRU

---

**Note**: This project is **open source**. Wayru, Inc and The Wayru Foundation are no longer operating entities, and will not provide any kind of support. The community is welcome to use, modify, and improve this codebase.
