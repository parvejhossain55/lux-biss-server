# Luxbiss Server

A production-ready, modular Go API server built with Gin Gonic, Postgres, and Redis. This project implements a robust authentication system and follows clean architecture principles.

## 🚀 Features

- **Authentication System**:
  - Email Registration & Login.
  - JWT Access & Refresh Token pairs.
  - Google OAuth Integration.
  - Password Recovery via Email OTP.
  - Secure Password Hashing (Bcrypt).
- **User Management**:
  - User Profiles.
  - Role-Based Access Control (RBAC) - `User` & `Admin`.
- **Performance & Security**:
  - Redis for OTP storage and caching.
  - Global Rate Limiting.
  - CORS and Security Headers middleware.
  - Structured Logging with Zap.
- **Developer Experience**:
  - Hot-reload with Air.
  - Full Docker & Docker Compose support.
  - Makefile for common tasks.
  - Integrated Database Migrations.

## 🛠️ Tech Stack

- **Lanuage**: Go (1.20+)
- **Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [GORM](https://gorm.io/)
- **Cache**: [Redis](https://redis.io/)
- **Auth**: JWT (RS256/HS256)
- **Validation**: Go Playground Validator v10
- **Configuration**: Viper

## 📁 Project Structure

```text
├── cmd/api             # Application entry point
├── docs/               # API documentation
├── internal/
│   ├── common/         # Shared utilities (Response, Pagination, etc.)
│   ├── config/         # Configuration logic
│   ├── database/       # Database connection drivers
│   ├── logger/         # Logging service
│   ├── middleware/     # Gin middlewares (Auth, RBAC, RateLimit)
│   ├── modules/        # Domain logic (Auth, User, Health)
│   └── server/         # Server initialization
├── migrations/         # SQL migration files
├── pkg/                # Reusable packages (JWT, Email, Hash)
└── Dockerfile          # Container definition
```

## ⚙️ Setup & Installation

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.20+ (if running locally)
- [Air](https://github.com/cosmtrek/air) (for hot-reload)

### 2. Environment Variables
Copy `.env.example` to `.env` and fill in your credentials:
```bash
cp .env.example .env
```

### 3. Running with Docker (easiest)
```bash
make docker-up
```
This starts the Go API, PostgreSQL, and Redis.

### 4. Running Locally
```bash
# Start Postgres & Redis only
docker compose up -d db redis

# Run the app with hot-reload
make dev
```

## 📖 API Documentation

The complete API reference with request/response examples can be found at:
👉 **[docs/API_DOCS.md](docs/API_DOCS.md)**

## 🧪 Testing & Quality

```bash
# Run tests
make test

# Check linting
make lint

# Tidy modules
make tidy
```

## 📜 License
Distributed under the MIT License. See `LICENSE` for more information.
