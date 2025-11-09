# Geekery API

<div align="center">

**Your Personal Geek Media Tracker**

A modern REST API for tracking and organizing geek media content (anime, movies, series, games, manga, light novels, music, books).

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 🎯 About

**Geekery** is a personal media tracking application inspired by MyAnimeList, IMDb, and Letterboxd. Built with clean architecture principles using Go, Gin, GORM, and PostgreSQL.

**Key Features:**
- 🔐 JWT Authentication with secure password hashing
- 📚 Global catalog of media items (shared across users)
- 📝 Personal lists with ratings, notes, and progress tracking
- 🏷️ Flexible tagging system
- 📊 Multi-completion support (re-watch/re-read tracking)

---

## 🚀 Quick Start

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (optional)

### Installation

```bash
# Clone repository
git clone https://github.com/rafaelc-rb/geekery-api.git
cd geekery-api

# Setup environment
cp .env.example .env
# Edit .env and set JWT_SECRET (minimum 32 characters)

# Start PostgreSQL + Run API
make dev

# Or step by step:
make docker-up     # Start PostgreSQL
make run           # Run API
```

The API will be available at `http://localhost:8080`

**Swagger Documentation:** `http://localhost:8080/swagger/index.html`

---

## 📖 API Overview

### Authentication

```bash
# Register
POST /api/auth/register
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepass123",
  "name": "John Doe"
}

# Login (returns JWT token)
POST /api/auth/login
{
  "username": "johndoe",  // or use email: "user@example.com"
  "password": "securepass123"
}
```

### Main Endpoints

- **Items (Catalog)**: `/api/items` - Global media catalog (public)
- **My List**: `/api/my-list` - Personal tracking (protected, requires JWT)
- **Tags**: `/api/tags` - Tag management
- **Health**: `/api/health` - Health check

**Protected routes require JWT token:**
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/my-list
```

📘 **Full API Documentation:** [Swagger UI](http://localhost:8080/swagger/index.html)

---

## 🏗 Architecture

```
├── cmd/                    # Application entrypoints
├── internal/
│   ├── auth/              # JWT & middleware
│   ├── config/            # Configuration
│   ├── database/          # DB connection & migrations
│   ├── dto/               # Data transfer objects
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Domain models
│   ├── repositories/      # Data access layer
│   ├── services/          # Business logic
│   └── routes/            # Route configuration
├── deploy/                # Docker compose
└── docs/                  # Swagger docs
```

**Clean Architecture:** Handlers → Services → Repositories → Database

---

## 🛠 Development

### Environment Variables

| Variable      | Description       | Required         |
| ------------- | ----------------- | ---------------- |
| `DB_HOST`     | PostgreSQL host   | ✅                |
| `DB_PORT`     | PostgreSQL port   | ✅                |
| `DB_USER`     | Database user     | ✅                |
| `DB_PASSWORD` | Database password | ✅                |
| `DB_NAME`     | Database name     | ✅                |
| `JWT_SECRET`  | JWT signing key   | ✅ (min 32 chars) |
| `SERVER_PORT` | API server port   | ✅                |
| `ENV`         | Environment       | ✅                |

### Make Commands

```bash
# Development
make dev           # Start everything (PostgreSQL + API)
make run           # Run API only
make seed          # Seed database with sample data

# Database
make docker-up     # Start PostgreSQL
make docker-down   # Stop PostgreSQL
make db-reset      # Reset database

# Testing & Quality
make test          # Run all tests
make test-unit     # Run unit tests only
make lint          # Run linter
make fmt           # Format code

# Build
make build         # Build binary
make clean         # Clean artifacts
```

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific test suites
make test-unit     # Unit tests (mocked dependencies)
make test-int      # Integration tests (real PostgreSQL)

# With coverage
go test -cover ./...
```

**Test Structure:**
- Unit tests: Services with mocked repositories
- Integration tests: Repository tests with real PostgreSQL
- Handler tests: HTTP handlers with mocked services

---

## 🛠 Tech Stack

| Technology     | Version | Purpose          |
| -------------- | ------- | ---------------- |
| **Go**         | 1.23+   | Backend language |
| **Gin**        | 1.11.0  | HTTP framework   |
| **GORM**       | 1.31.1  | ORM              |
| **PostgreSQL** | 18      | Database         |
| **JWT**        | v5.3.0  | Authentication   |
| **bcrypt**     | -       | Password hashing |
| **Docker**     | Latest  | Containerization |

---

## 🗺 Roadmap

### ✅ Phase 1: MVP (Completed)
- Clean architecture setup
- CRUD operations for items and user lists
- JWT authentication
- Flexible progress tracking
- Comprehensive testing

### 🎯 Phase 2: Performance & Scale (Next)
- Database indices optimization
- Query optimization (N+1 fixes)
- Pagination for all endpoints
- Full-text search
- Rate limiting

### 📅 Phase 3: Advanced Features
- External API integrations (MAL, TMDB, IGDB)
- Image uploads (S3/CloudFlare)
- Social features
- Recommendations engine
- Advanced statistics

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Rafael C Ribeiro**

- GitHub: [@rafaelc-rb](https://github.com/rafaelc-rb)

---

<div align="center">

**Made with ❤️ for geeks, by geeks**

⭐ Star this repo if you find it helpful!

</div>
