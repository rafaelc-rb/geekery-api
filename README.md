# Geekery API

<div align="center">

**Your Personal Geek Media Tracker**

A modern REST API for tracking and organizing geek media content (anime, movies, series, games, manga, light novels, music, books).

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 📋 Table of Contents

- [About](#about)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Roadmap](#roadmap)

---

## 🎯 About

**Geekery** is a personal media tracking application inspired by MyAnimeList, IMDb, and Letterboxd. This repository contains the backend API built with clean architecture principles, designed to be scalable and maintainable.

### Key Concepts

- **Items**: Global catalog of media (anime, movies, games, etc.) - shared across all users
- **User Lists**: Personal tracking with status, ratings, notes, and progress
- **Flexible Progress**: Track episodes, chapters, playtime, or completion percentage
- **Multi-view Support**: Re-watch/re-read tracking with history

---

## ✨ Features

### Current (MVP)

- ✅ **Global Catalog Management**: CRUD operations for media items
- ✅ **Personal Lists**: Track your media with status, ratings, and notes
- ✅ **Flexible Progress Tracking**:
  - Episodic (anime, series)
  - Reading (manga, books)
  - Time-based (movies)
  - Percentage (games)
  - Play count (music)
- ✅ **Tags System**: Many-to-many relationships for categorization
- ✅ **Multi-completion Support**: Track re-watches and re-reads
- ✅ **Type-specific Data**: JSONB storage for specialized metadata
- ✅ **Clean Architecture**: Handlers → Services → Repositories
- ✅ **PostgreSQL 18**: Latest stable database with GORM ORM
- ✅ **Comprehensive Testing**: Unit tests + Integration tests
- ✅ **Docker Ready**: Full containerization support

### Coming Soon

- 🔜 **JWT Authentication**: User registration and login
- 🔜 **Pagination**: Efficient data loading for large lists
- 🔜 **Advanced Search**: Full-text search with filters
- 🔜 **Image Uploads**: Cover images for items
- 🔜 **External APIs**: Integration with MAL, TMDB, IGDB
- 🔜 **Statistics Dashboard**: Analytics and insights

---

## 🛠 Tech Stack

| Technology     | Version | Purpose                     |
| -------------- | ------- | --------------------------- |
| **Go**         | 1.23+   | Backend language            |
| **Gin**        | 1.11.0  | HTTP framework              |
| **GORM**       | 1.31.1  | ORM for database operations |
| **PostgreSQL** | 18      | Primary database            |
| **Docker**     | Latest  | Containerization            |
| **godotenv**   | 1.5.1   | Environment configuration   |

---

## 🏗 Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│          HTTP Layer (Gin)               │
│  ┌─────────────────────────────────┐   │
│  │     Handlers / Controllers      │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│         Business Logic Layer            │
│  ┌─────────────────────────────────┐   │
│  │          Services               │   │
│  │  (Validation, Business Rules)   │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│          Data Access Layer              │
│  ┌─────────────────────────────────┐   │
│  │        Repositories             │   │
│  │    (GORM Database Access)       │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│            PostgreSQL 18                │
└─────────────────────────────────────────┘
```

### Project Structure

```
geekery-api/
├── cmd/
│   ├── main.go                    # Application entry point
│   └── seed/
│       └── main.go                # Database seeder
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── database/
│   │   └── database.go            # Database connection & migrations
│   ├── models/
│   │   ├── item.go                # Global catalog item
│   │   ├── user_item.go           # Personal list entry
│   │   ├── tag.go                 # Tag model
│   │   ├── types.go               # Enums and types
│   │   ├── user_item_helpers.go   # Progress tracking helpers
│   │   ├── *_data.go              # Type-specific models
│   │   └── errors.go              # Custom errors
│   ├── dto/
│   │   ├── item_dto.go            # Item DTOs
│   │   ├── user_item_dto.go       # UserItem DTOs
│   │   ├── mappers.go             # Model to DTO mappers
│   │   └── ...
│   ├── repositories/
│   │   ├── interfaces.go          # Repository interfaces
│   │   ├── item_repository.go     # Catalog data access
│   │   ├── user_item_repository.go# User list data access
│   │   └── tag_repository.go      # Tag data access
│   ├── services/
│   │   ├── item_service.go        # Catalog business logic
│   │   ├── user_item_service.go   # User list business logic
│   │   └── tag_service.go         # Tag business logic
│   ├── handlers/
│   │   ├── item_handler.go        # Catalog HTTP handlers
│   │   ├── user_item_handler.go   # User list HTTP handlers
│   │   └── tag_handler.go         # Tag HTTP handlers
│   ├── routes/
│   │   └── routes.go              # Route configuration
│   ├── middleware/
│   │   └── logger.go              # Request logging
│   ├── logger/
│   │   └── logger.go              # Structured logging
│   └── testutil/
│       ├── mocks.go               # Mock repositories
│       ├── fixtures.go            # Test fixtures
│       └── db.go                  # Test database setup
├── test/
│   └── integration/
│       └── repository_integration_test.go
├── deploy/
│   └── docker-compose.yml         # PostgreSQL container
├── .env.example                   # Environment template
├── .gitignore                     # Git ignore rules
├── Makefile                       # Development commands
├── go.mod                         # Go dependencies
└── README.md                      # This file
```

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (optional, but recommended)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/rafaelc-rb/geekery-api.git
cd geekery-api
```

2. **Set up environment variables**

```bash
cp .env.example .env
# Edit .env if needed (defaults work for development)
```

3. **Quick setup (recommended)**

```bash
make setup
# This will: install dependencies, start PostgreSQL, create database
```

4. **Run the application**

```bash
make run
```

**Alternative: One-command setup**

```bash
make dev
# Starts PostgreSQL, configures database, and runs the API
```

The API will be available at `http://localhost:8080`

### Seed Database (Optional)

```bash
make seed
# Populates database with sample data (anime, movies, games, etc.)
```

---

## 📖 API Documentation

### Base URL

```
http://localhost:8080/api
```

### Endpoints

#### Health Check

```http
GET /api/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Geekery API is running!"
}
```

---

### 📚 Catalog Endpoints (Items)

Global catalog shared across users (admin-only in production).

#### Create Item

```http
POST /api/items
Content-Type: application/json

{
  "title": "Attack on Titan",
  "type": "anime",
  "description": "Humanity fights titans",
  "release_date": "2013-04-07T00:00:00Z",
  "cover_url": "https://example.com/cover.jpg",
  "tag_ids": [1, 2, 3]
}
```

**Valid Types:** `anime`, `movie`, `series`, `game`, `manga`, `light_novel`, `music`, `book`

#### Get All Items

```http
GET /api/items
GET /api/items?type=anime  # Filter by type
```

#### Search Items

```http
GET /api/items/search?q=attack
```

#### Get Item by ID

```http
GET /api/items/:id
```

#### Update Item

```http
PUT /api/items/:id
Content-Type: application/json

{
  "title": "Attack on Titan Final Season",
  "description": "The final arc"
}
```

#### Delete Item

```http
DELETE /api/items/:id
```

---

### 🗂 Personal List Endpoints (My List)

Manage your personal tracking list.

#### Add Item to My List

```http
POST /api/my-list
Content-Type: application/json

{
  "item_id": 1,
  "status": "in_progress"
}
```

**Valid Statuses:** `planned`, `in_progress`, `completed`, `paused`, `dropped`

#### Get My List

```http
GET /api/my-list
GET /api/my-list?status=completed    # Filter by status
GET /api/my-list?favorite=true       # Filter favorites
```

#### Get List Item Details

```http
GET /api/my-list/:id
```

#### Update List Item

```http
PUT /api/my-list/:id
Content-Type: application/json

{
  "status": "completed",
  "rating": 9.5,
  "favorite": true,
  "notes": "Amazing anime!",
  "progress_type": "episodic",
  "progress_data": {
    "season": 4,
    "episode": 28
  }
}
```

**Progress Types:**
- `episodic`: For anime/series (season, episode)
- `reading`: For books/manga (chapter, volume, page)
- `time`: For movies (minutes_watched)
- `percent`: For games (percent, hours)
- `boolean`: For music (listened, play_count)

#### Remove from List

```http
DELETE /api/my-list/:id
```

#### Get Statistics

```http
GET /api/my-list/stats
```

**Response:**
```json
{
  "total": 150,
  "in_progress": 12,
  "completed": 95,
  "planned": 38,
  "paused": 3,
  "dropped": 2,
  "favorites": 25
}
```

---

### 🏷 Tags Endpoints

#### Create Tag

```http
POST /api/tags
Content-Type: application/json

{
  "name": "Action"
}
```

#### Get All Tags

```http
GET /api/tags
```

#### Get Tag by ID

```http
GET /api/tags/:id
```

#### Update Tag

```http
PUT /api/tags/:id
Content-Type: application/json

{
  "name": "Adventure"
}
```

#### Delete Tag

```http
DELETE /api/tags/:id
```

---

### Example Usage with cURL

**Complete workflow:**

```bash
# 1. Create tags
curl -X POST http://localhost:8080/api/tags \
  -H "Content-Type: application/json" \
  -d '{"name":"Action"}'

curl -X POST http://localhost:8080/api/tags \
  -H "Content-Type: application/json" \
  -d '{"name":"Shonen"}'

# 2. Create an anime in catalog
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Attack on Titan",
    "type": "anime",
    "description": "Humanity fights against titans",
    "release_date": "2013-04-07T00:00:00Z",
    "tag_ids": [1, 2]
  }'

# 3. Add to your personal list
curl -X POST http://localhost:8080/api/my-list \
  -H "Content-Type: application/json" \
  -d '{
    "item_id": 1,
    "status": "in_progress"
  }'

# 4. Update progress
curl -X PUT http://localhost:8080/api/my-list/1 \
  -H "Content-Type: application/json" \
  -d '{
    "progress_type": "episodic",
    "progress_data": {
      "season": 4,
      "episode": 28
    }
  }'

# 5. Mark as completed with rating
curl -X PUT http://localhost:8080/api/my-list/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "rating": 9.5,
    "favorite": true,
    "notes": "One of the best anime ever!"
  }'

# 6. View statistics
curl http://localhost:8080/api/my-list/stats
```

---

## 💻 Development

### Available Make Commands

```bash
# Quick Start
make setup         # Complete initial setup (first time only)
make dev           # Start everything (PostgreSQL + API)
make run           # Run the API (PostgreSQL must be running)
make seed          # Seed database with sample data

# Database Commands
make docker-up     # Start PostgreSQL container
make docker-down   # Stop PostgreSQL container
make docker-logs   # View PostgreSQL logs
make db-reset      # Reset database (removes all data)
make db-verify     # Test database connection

# Development
make build         # Build application binary
make test          # Run all tests
make test-unit     # Run unit tests only
make test-int      # Run integration tests only
make clean         # Clean build artifacts
make fmt           # Format Go code
make lint          # Run linters
make tidy          # Tidy Go dependencies

# Help
make help          # Show all available commands
```

### Environment Variables

| Variable      | Description       | Default                     |
| ------------- | ----------------- | --------------------------- |
| `DB_HOST`     | PostgreSQL host   | `localhost`                 |
| `DB_PORT`     | PostgreSQL port   | `5433`                      |
| `DB_USER`     | Database user     | `geekery`                   |
| `DB_PASSWORD` | Database password | `your_secure_password_here` |
| `DB_NAME`     | Database name     | `geekery_db`                |
| `SERVER_PORT` | API server port   | `8080`                      |
| `ENV`         | Environment       | `development`               |
| `LOG_LEVEL`   | Logging level     | `info`                      |

---

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-int

# With coverage
go test -cover ./...
```

### Test Structure

- **Unit Tests**: Services with mocked repositories
- **Integration Tests**: Repository tests with real PostgreSQL
- **Handler Tests**: HTTP handlers with mocked services

---

## 🗺 Roadmap

### Phase 1: MVP ✅ (COMPLETED)

- [x] Clean architecture setup
- [x] Global catalog (Items) CRUD
- [x] Personal lists (UserItems) CRUD
- [x] Tags with many-to-many relationships
- [x] Flexible progress tracking system
- [x] Multi-completion support (re-watch/re-read)
- [x] Type-specific data models
- [x] Comprehensive testing
- [x] Docker setup
- [x] Database seeder

### Phase 2: Authentication 🎯 (NEXT)

- [ ] JWT implementation
- [ ] User registration & login
- [ ] Password hashing with bcrypt
- [ ] Protected routes middleware
- [ ] User-specific data isolation
- [ ] Refresh token support

### Phase 3: Performance & Scale

- [ ] Pagination for all list endpoints
- [ ] Full-text search with PostgreSQL
- [ ] Caching layer (Redis)
- [ ] Database indices optimization
- [ ] Query optimization (N+1 fixes)
- [ ] Rate limiting

### Phase 4: Advanced Features

- [ ] Image uploads with S3/CloudFlare
- [ ] External API integrations (MAL, TMDB, IGDB)
- [ ] Reviews and comments system
- [ ] Social features (follow users, activity feed)
- [ ] Recommendations engine
- [ ] Import/export functionality

### Phase 5: Production Ready

- [ ] OpenAPI/Swagger documentation
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Monitoring and observability
- [ ] Error tracking (Sentry)
- [ ] Performance metrics
- [ ] Backup strategies
- [ ] Deployment guides (AWS, GCP, etc.)

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Rafael C Ribeiro**

- GitHub: [@rafaelc-rb](https://github.com/rafaelc-rb)

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

Feel free to check the [issues page](https://github.com/rafaelc-rb/geekery-api/issues).

### How to Contribute

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'feat: add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

<div align="center">

**Made with ❤️ for geeks, by geeks**

⭐ Star this repo if you find it helpful!

</div>
