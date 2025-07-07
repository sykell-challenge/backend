# Sykell Challenge - Backend API

A Go-based REST API for the URL Management System. This backend provides authentication, URL management, and analytics functionality.

## Features

### URL Management
- Create, read, update, and delete URLs
- Track URL processing status (queued, running, done, error)
- HTML version detection
- Login form detection
- Tag system for URL categorization
- Link analysis (internal, external, broken links)

### Search & Analytics
- Basic URL search
- Fuzzy search functionality
- URL statistics and analytics

### User Management
- User registration and authentication
- JWT-based authorization
- Protected routes for authenticated users

## API Endpoints

### Public Routes
- `POST /users` - User registration
- `POST /users/login` - User login

### Protected Routes (require JWT authentication)

**URL Management:**
- `GET /urls` - List all URLs
- `GET /urls/search` - Search URLs by string
- `GET /urls/search/fuzzy` - Fuzzy search URLs
- `GET /urls/stats` - Get URL statistics
- `GET /urls/:id` - Get specific URL
- `GET /urls/:id/links` - Get all links from URL
- `GET /urls/:id/links/internal` - Get internal links
- `GET /urls/:id/links/external` - Get external links
- `GET /urls/:id/links/broken` - Get broken links
- `POST /urls` - Create new URL
- `PUT /urls/:id` - Update URL
- `PATCH /urls/:id/status` - Update URL status
- `DELETE /urls/:id` - Delete URL

**User Management:**
- `GET /users` - List all users
- `GET /users/:id` - Get specific user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

## Technology Stack

- **Go** 1.21+
- **Gin** Web Framework
- **GORM** (Go ORM)
- **JWT** for authentication
- **MySQL** driver
- **bcrypt** for password hashing

## Project Structure

```
backend/
├── auth/                 # JWT authentication
│   ├── jwt.go
│   └── middleware.go
├── db/                   # Database configuration
│   └── db.go
├── handlers/             # HTTP handlers
│   ├── url/             # URL-related handlers
│   │   ├── handler.go
│   │   ├── url_crud.go
│   │   ├── url_links.go
│   │   ├── url_list.go
│   │   └── url_search.go
│   └── user/            # User-related handlers
│       ├── handler.go
│       ├── user_auth.go
│       ├── user_crud.go
│       └── user_list.go
├── models/              # Data models
│   ├── links.go
│   ├── tags.go
│   ├── urls.go
│   └── user.go
├── repositories/        # Data access layer
│   ├── link_repository.go
│   ├── tag_repository.go
│   ├── url_repository.go
│   └── user_repository.go
├── main.go              # Application entry point
├── go.mod
└── go.sum
```

## Environment Setup

### Prerequisites

- Go 1.21+
- MySQL 8.0+
- Git

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <backend-repo-url>
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file or set the following environment variables:
   ```bash
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=sykell
   DB_PASSWORD=sykellpass
   DB_NAME=websites_dev
   ALLOWED_ORIGINS=http://localhost:3000
   JWT_SECRET=your-secret-key
   ```

4. **Run database migrations**
   The application will automatically create the required tables on startup.

5. **Start the development server**
   ```bash
   go run main.go
   ```

   The API will be available at http://localhost:8080

### Docker Development

1. **Build the Docker image**
   ```bash
   docker build -f Dockerfile.dev -t sykell-backend-dev .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

### Production Deployment

1. **Build the production image**
   ```bash
   docker build -f Dockerfile -t sykell-backend .
   ```

2. **Run the production container**
   ```bash
   docker run -p 8080:8080 \
     -e DB_HOST=your-db-host \
     -e DB_USER=your-db-user \
     -e DB_PASSWORD=your-db-password \
     -e DB_NAME=your-db-name \
     -e ALLOWED_ORIGINS=https://your-frontend-domain.com \
     sykell-backend
   ```

## Environment Variables

- `DB_HOST` - Database host (default: `localhost`)
- `DB_PORT` - Database port (default: `3306`)
- `DB_USER` - Database username (required)
- `DB_PASSWORD` - Database password (required)
- `DB_NAME` - Database name (required)
- `ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins
- `JWT_SECRET` - Secret key for JWT token signing (required)
- `PORT` - Server port (default: `8080`)

## Database Schema

The application uses the following main models:

### Users
- ID, Username, Email, PasswordHash
- FirstName, LastName, IsActive
- CreatedAt, UpdatedAt, LastLoginAt

### URLs
- ID, URL, Status, HTMLVersion
- LoginForm (boolean), Tags (JSON), Links (JSON)
- CreatedAt, UpdatedAt

### Tags
- Custom JSON structure for URL categorization

### Links
- Custom JSON structure for internal, external, and broken links

## API Testing

Use the included `postman_collection.json` for testing the API endpoints, or test with curl:

```bash
# Register a new user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","first_name":"Test","last_name":"User"}'

# Login
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Create a URL (requires JWT token)
curl -X POST http://localhost:8080/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"url":"https://example.com","tags":{}}'
```

## Development

### Code Organization

The backend follows a clean architecture pattern:

- **Handlers**: HTTP request/response handling
- **Repositories**: Data access layer
- **Models**: Data structures and validation
- **Auth**: Authentication and authorization middleware

### Adding New Features

1. Create model in `models/`
2. Create repository in `repositories/`
3. Create handlers in `handlers/`
4. Register routes in `main.go`

### Running Tests

```bash
go test ./...
```

## License

This project is part of the Sykell technical challenge.
