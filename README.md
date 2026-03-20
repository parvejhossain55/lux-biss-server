# Luxbiss Server

Luxbiss Server is a modular Go API built with Gin, GORM, PostgreSQL, and Redis. It provides authentication, user and product management, wallets, transactions, gift cards, file uploads, admin reporting, and embedded frontend asset serving.

## Features

- JWT-based authentication with access and refresh tokens
- Google OAuth login
- Password reset flow with email OTP verification
- Role-based access control for `admin` and `user`
- Product, level, and step management
- Wallet, transaction, and gift card workflows
- Manager management and admin dashboard stats
- Image upload endpoint with static file serving from `/uploads`
- Redis-backed rate limiting and OTP storage
- Structured logging, CORS, CSRF, recovery, and security middleware
- Database migrations and seeders via Make targets
- Embedded frontend assets served by the API when available

## Tech Stack

- Go `1.25.0`
- Gin
- GORM
- PostgreSQL
- Redis
- Viper
- Zap
- Google OAuth API client

## Project Structure

```text
.
├── cmd/
│   ├── api/                 # API entrypoint and embedded frontend serving
│   └── seeder/              # Seeder entrypoint
├── docs/                    # Reserved for project docs
├── internal/
│   ├── common/              # Shared helpers and response utilities
│   ├── config/              # Environment-based configuration loading
│   ├── database/            # Postgres and Redis initialization
│   ├── logger/              # Application logger setup
│   ├── middleware/          # Auth, RBAC, CORS, CSRF, rate limit, recovery
│   ├── modules/             # Feature modules
│   │   ├── admindashboard/
│   │   ├── auth/
│   │   ├── giftcard/
│   │   ├── health/
│   │   ├── manager/
│   │   ├── product/
│   │   ├── transaction/
│   │   ├── upload/
│   │   ├── user/
│   │   └── wallet/
│   └── server/              # HTTP server lifecycle
├── migrations/              # SQL migrations
├── pkg/                     # Reusable packages (email, hash, jwt, validators)
├── public/uploads/          # Uploaded files served at /uploads
├── scripts/                 # Helper scripts
├── .air.toml                # Air hot-reload config
├── .golangci.yml            # Linter config
├── Makefile
├── go.mod
└── go.sum
```

## Requirements

- Go `1.25.0`
- PostgreSQL
- Redis
- `air` for hot reload during development
- `golangci-lint` for linting
- `migrate` CLI if you want to run SQL migrations manually

## Configuration

The server loads configuration from a local `.env` file and environment variables.

There is currently no committed `.env.example`, so create `.env` manually in the project root. Important settings include:

```env
APP_NAME=luxbiss-server
APP_ENV=dev

SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=luxbiss
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=replace-with-a-secret-at-least-32-characters-long
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=168h
JWT_ISSUER=luxbiss-server

COOKIE_DOMAIN=localhost
COOKIE_SECURE=false
COOKIE_SAMESITE=lax
COOKIE_PATH=/

LOG_LEVEL=info
LOG_FORMAT=console

CORS_ALLOWED_ORIGINS=http://localhost:3000

SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=

BOT_TOKEN=
TELEGRAM_CHAT_ID=
TELEGRAM_PROXY=
```

Notes:

- `JWT_SECRET` is required and must be at least 32 characters.
- In production, `DB_SSLMODE=disable` is rejected by config validation.
- TLS can be enabled with `SERVER_TLS_ENABLED=true` plus certificate and key file paths.

## Running the Server

Install dependencies and start the API locally:

```bash
go mod tidy
make run
```

For hot reload during development:

```bash
make dev
```

The API starts on `http://localhost:8080` by default.

## Available Make Targets

```bash
make dev            # Run with Air hot reload
make run            # Run the API directly
make build          # Build the binary into ./bin
make test           # Run tests
make test-cover     # Generate coverage.out and coverage.html
make lint           # Run golangci-lint
make tidy           # Tidy and verify Go modules
make migrate-up     # Apply all SQL migrations
make migrate-down   # Roll back one migration
make migrate-create name=create_example
make seed           # Seed database data
make seed-fresh     # Truncate and reseed
```

## API Base Path

All API routes are registered under:

```text
/api/v1
```

Current route groups include:

- `/auth`
- `/health`
- `/users`
- `/products`
- `/levels`
- `/wallets`
- `/managers`
- `/transactions`
- `/giftcards`
- `/admin/dashboard`
- `/upload`

Uploads are served statically from:

```text
/uploads
```

## Notes on Frontend Assets

If static frontend files are present under `cmd/api/web/out`, the server embeds and serves them for non-API routes. If they are absent during development, the API still runs normally.