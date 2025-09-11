# Ghostvox Backend RESTful API

A production-intent RESTful API written in Go that powers the Ghostvox application. It exposes endpoints for managing polls and related domain data. The service is containerized, deployable on Fly.io, and persists data in PostgreSQL.

## Table of Contents
- [Key Features](#key-features)
- [Architecture Overview](#architecture-overview)
- [OpenAPI Specification](#openapi-specification)
- [Authentication & Authorization](#authentication--authorization)
- [Error Handling](#error-handling)
- [Local Development](#local-development)
- [Local HTTPS & OAuth](#local-https--oauth)
- [Environment Variables](#environment-variables)
- [Running with Docker](#running-with-docker)
- [Deployment (Fly.io)](#deployment-flyio)
- [Database & Data Integrity](#database--data-integrity)
- [Example Requests](#example-requests)
- [Project Structure (High-Level)](#project-structure-high-level)
- [Contributing](#contributing)
- [License](#license)

## Key Features
- Go-based RESTful API
- PostgreSQL persistence
- JWT-based auth (access + refresh tokens)
- Secure HTTP-only cookie storage for auth tokens
- Role-Based Access Control (admin-only endpoints)
- Google & GitHub OAuth integration
- OpenAPI spec (`openapi_spec.yaml`) as the single source of truth
- Transactional integrity for critical operations (e.g. user creation, token workflows)
- Consistent, machine- & human-friendly error format
- Docker & Docker Compose for reproducible environments
- Designed for Fly.io deployment

## Architecture Overview
The service follows conventional Go layering:

1. Transport / HTTP handlers
2. Service / business logic
3. Data access (repositories) backed by PostgreSQL
4. Cross-cutting: auth, validation, config, logging

The OpenAPI document drives:
- Contract-first design
- Documentation & interactive exploration
- Potential future code generation (clients / mocks)

## OpenAPI Specification
The source of truth for all endpoints is `openapi_spec.yaml` located at the repository root.

### Viewing the Specification
1. Open `openapi_spec.yaml`
2. Copy all contents
3. Visit https://editor.swagger.io/
4. Paste into the editor for an interactive UI

Avoid documenting endpoints manually in this README to prevent drift—update the spec instead.

## Authentication & Authorization
- Authentication uses short-lived Access Tokens and longer-lived Refresh Tokens
- Both tokens are issued as signed JWTs
- Tokens are stored in secure HTTP-only cookies to mitigate XSS access
- Admin-only endpoints enforce role checks
- OAuth providers (Google, GitHub) supported (requires HTTPS redirect URIs)

(If you need exact cookie names, expiration durations, or claim structure, inspect the implementation or OpenAPI schema—they are intentionally not duplicated here to keep a single source of truth.)

### Typical Flow
1. User authenticates (credentials or OAuth)
2. Server sets access + refresh token cookies
3. Client calls protected endpoints (access token)
4. On expiry, client hits refresh endpoint (refresh token)
5. If refresh valid, new token pair issued
6. Logout / revoke clears cookies / invalidates refresh lineage

## Error Handling
All error responses follow a consistent structured format to simplify client-side validation & UX.

### JSON Structure
```json
{
  "errors": {
    "field_name": "Error message specific to this field",
    "another_field": "Another validation or domain error"
  }
}
```

### HTTP Status Codes (Representative)
- 400 Bad Request: Invalid input / validation failure
- 401 Unauthorized: Missing or invalid authentication
- 403 Forbidden: Authenticated but insufficient privileges
- 404 Not Found: Resource does not exist
- 409 Conflict: Resource already exists (e.g. duplicate email)
- 500 Internal Server Error: Unhandled server condition

This predictable contract enables straightforward form + inline error displays on the client.

## Local Development
Choose either native Go tooling or Docker. Docker Compose is recommended for parity.

### Prerequisites
- Go (if building locally)
- Docker & Docker Compose
- PostgreSQL (if not using the Compose service)
- OpenSSL / mkcert (for generating local certificates)

### Quick Start (Docker Compose)
```bash
docker compose up --build
```

By default the API will be reachable at:
- http://localhost:8080 (if HTTP)
- https://localhost:8080 (if HTTPS enabled)

### Live Reload (If configured)
If an auto-reload tool is integrated (e.g. `air`), changes will rebuild automatically. Otherwise restart the container.

## Local HTTPS & OAuth
HTTPS is required for local OAuth (Google / GitHub) redirect flows.

Steps:
1. Generate or obtain development certificates named:
   - `localhost+2.pem`
   - `localhost+2-key.pem`
2. Place both files in the project root
3. Set `USE_HTTPS=true` in `.env`
4. Ensure OAuth app redirect URIs use https://localhost:8080

(Certificate generation can be done with mkcert: `mkcert localhost 127.0.0.1 ::1`)

## Environment Variables
Full list of currently expected environment variables (correct spellings shown; fix any typos in your local `.env`):

| Variable | Purpose / Description |
|----------|-----------------------|
| MODE | Runtime mode (e.g. development, production) - may toggle logging / security defaults |
| ACCESS_ORIGIN | Allowed origin (CORS) for browser clients (e.g. https://app.example.com) |
| ACCESS_TOKEN_EXPIRES | Access token lifetime (e.g. 15m, 1h) used when issuing JWTs |
| USE_HTTPS | Set to true to enable TLS locally |
| TLS_CERT_PATH | Filesystem path to TLS certificate (used when USE_HTTPS=true) |
| TLS_KEY_PATH | Filesystem path to TLS private key (used when USE_HTTPS=true) |
| GOOGLE_CLIENT_ID | Google OAuth Client ID |
| GOOGLE_CLIENT_SECRET | Google OAuth Client Secret |
| GOOGLE_REDIRECT_URI | Google OAuth redirect URI (must match console) |
| GITHUB_CLIENT_ID | GitHub OAuth Client ID |
| GITHUB_CLIENT_SECRET | GitHub OAuth Client Secret |
| GITHUB_REDIRECT_URI | GitHub OAuth redirect URI |
| CRON_CHECK_FOR_EXPIRED_POLLS | Cron expression / interval controlling scheduled poll expiration task |
| AWS_ACCESS_KEY_ID | AWS credential for S3 access |
| AWS_SECRET_ACCESS_KEY | AWS secret credential for S3 access |
| AWS_REGION | AWS region of the S3 bucket |
| AWS_S3_BUCKET | S3 bucket name used for asset/object storage |
| IP_RATE_LIMIT | Base allowed requests per interval (token refill rate) |
| IP_RATE_BURST | Burst capacity above steady rate (token bucket size) |

Keep secrets out of version control—use a local `.env` or managed secret store in production.

## Running with Docker

### Build Image
```bash
docker build -t ghostvox-backend .
```

### Run Container
```bash
docker run --rm -p 8080:8080 ghostvox-backend
```

### Compose (App + DB)
See the earlier Quick Start; typically Compose defines both API + PostgreSQL.

## Deployment (Fly.io)
High-level deployment outline (verify with your Fly configuration files):

1. Authenticate with Fly
2. Create / configure app (`fly launch` or existing `fly.toml`)
3. Provision PostgreSQL (Fly Postgres cluster)
4. Set secrets (`fly secrets set KEY=VALUE`)
5. Deploy (`fly deploy`)
6. Run any migrations (depending on tooling)

Monitor logs and scaling using Fly's CLI or dashboard.

## Database & Data Integrity
- PostgreSQL with foreign keys & cascading deletes to maintain referential integrity
- Transactions wrap multi-step critical operations (e.g., new user + issuing tokens)
- Poll-related entities likely reference users (consult schema for relations)

(Refer to migration or model files for exact table definitions; not duplicated here.)

## Example Requests

### Health / Root (Example)
```bash
curl -i http://localhost:8080/health
```

### Authenticated Request (Pseudo)
```bash
curl -i \
  -H "Cookie: access_token=YOUR_ACCESS_JWT" \
  http://localhost:8080/api/protected-resource
```

### Refresh Token
```bash
curl -i \
  -H "Cookie: refresh_token=YOUR_REFRESH_JWT" \
  -X POST http://localhost:8080/auth/refresh
```

(Refer to OpenAPI spec for exact paths, verbs, and schemas.)

## Project Structure (High-Level)
(Actual paths may differ—adjust to match the repository.)

```
/cmd            Main entrypoint(s)
/internal       Application code (services, handlers, repos, auth)
/pkg            Reusable packages (if any)
openapi_spec.yaml
docker-compose.yml
Dockerfile
```

## License
[MIT](LICENSE.txt)

---

For any gaps or ambiguities, prefer inspecting the code and OpenAPI document—this README intentionally avoids duplication to reduce drift.
