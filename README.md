# Chirpy - Social Media API

A Twitter-like social media API built in Go, featuring user authentication, post management, and premium subscription handling through webhook integrations.

## Overview

Chirpy is a RESTful API that demonstrates modern backend development practices in Go. It provides a complete social media platform backend with user management, content creation, and payment integration capabilities.

**Key Features:**
- 🔐 **JWT Authentication** - Secure user authentication with access and refresh tokens
- 📝 **Content Management** - Create, read, update, and delete chirps (posts)
- 👥 **User Management** - User registration, login, and profile updates
- 💎 **Premium Subscriptions** - Chirpy Red premium memberships via webhook integration
- 🗄️ **Database Migrations** - Structured PostgreSQL schema management
- 🔍 **Advanced Filtering** - Query chirps by author and sort by date
- 🔒 **Secure API Keys** - Protected webhook endpoints with API key authentication

## Why This Project Matters

This project showcases production-ready backend development skills including:

- **Clean Architecture**: Well-structured Go codebase with separation of concerns
- **Database Design**: PostgreSQL with proper relationships and constraints
- **Security Best Practices**: JWT tokens, password hashing, API key validation
- **API Design**: RESTful endpoints with proper HTTP status codes
- **Testing**: Comprehensive test suite for authentication and business logic
- **DevOps**: Database migrations and environment configuration

Perfect for demonstrating backend engineering capabilities to potential employers or as a foundation for larger projects.

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 14+
- **Authentication**: JWT with refresh tokens
- **Database Toolkit**: sqlc for type-safe SQL
- **Migrations**: goose
- **Password Hashing**: bcrypt
- **HTTP Router**: Go standard library

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- `sqlc` and `goose` CLI tools

```bash
# Install required tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd chirpy
   ```

2. **Set up PostgreSQL**
   ```bash
   # Create database
   sudo -u postgres createdb chirpy
   
   # Start PostgreSQL service
   sudo systemctl start postgresql
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and secrets
   ```

4. **Run database migrations**
   ```bash
   goose -dir sql/schema postgres "your-connection-string" up
   ```

5. **Generate database code**
   ```bash
   sqlc generate
   ```

6. **Build and run**
   ```bash
   go build -o chirpy
   ./chirpy
   ```

The server will start on `http://localhost:8080`

## API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/api/users` | Create user account | None |
| POST | `/api/login` | User login | None |
| POST | `/api/refresh` | Refresh access token | Refresh Token |
| POST | `/api/revoke` | Revoke refresh token | Refresh Token |
| PUT | `/api/users` | Update user profile | Access Token |

### Chirp Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/api/chirps` | Get all chirps | None |
| GET | `/api/chirps?author_id={id}` | Get chirps by author | None |
| GET | `/api/chirps?sort=desc` | Get chirps sorted by date | None |
| POST | `/api/chirps` | Create new chirp | Access Token |
| DELETE | `/api/chirps/{id}` | Delete chirp | Access Token |

### Webhook Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/api/polka/webhooks` | Handle payment webhooks | API Key |

### Admin Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| GET | `/admin/metrics` | Server metrics | None (dev only) |
| POST | `/admin/reset` | Reset database | None (dev only) |

### Example Requests

**Create User:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

**Create Chirp:**
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"body":"Hello, world!"}'
```

## Project Structure

```
chirpy/
├── main.go                 # Application entry point
├── types.go               # Type definitions
├── handlers_*.go          # HTTP handlers
├── internal/
│   ├── auth/             # Authentication logic
│   └── database/         # Generated database code
├── sql/
│   ├── schema/           # Database migrations
│   └── queries/          # SQL queries
└── assets/               # Static files
```

## Environment Configuration

Required environment variables:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-jwt-secret-here
POLKA_KEY=your-polka-api-key
```

## Development

### Running Tests
```bash
go test ./...
```

### Adding Database Changes
1. Create migration: `goose -dir sql/schema create migration_name sql`
2. Write up/down migrations
3. Run migration: `goose -dir sql/schema postgres "connection-string" up`
4. Regenerate code: `sqlc generate`

### Code Generation
- Database models and queries are generated using `sqlc`
- Run `sqlc generate` after any SQL changes

## Security Features

- **Password Security**: bcrypt hashing with salt
- **Token Security**: JWT with short expiration + refresh tokens
- **API Protection**: Bearer token authentication
- **Webhook Security**: API key validation for external services
- **SQL Injection Prevention**: Parameterized queries via sqlc

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is open source and available under the MIT License.