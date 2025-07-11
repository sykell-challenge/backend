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

## Environment Setup

### Prerequisites

- Go 1.21+
- MySQL 8.0+
- Git

### Development Setup

## License

This project is part of the Sykell technical challenge.
