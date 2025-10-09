# WorkoutPal-Backend

Backend API for the WorkoutPal gym workout tracking application built with Go.

## Project Overview

WorkoutPal is a comprehensive fitness tracking platform that helps users plan, log, and visualize their workout progress. This repository contains the backend API that powers the application.

## Features

- **User Management**: Create, read, update, delete user profiles
- **Goal Tracking**: Set and monitor fitness goals with deadlines
- **Social Features**: Follow other users and view their progress
- **Workout Routines**: Create and manage custom exercise routines
- **Database Support**: PostgreSQL with fallback to in-memory storage
- **REST API**: Clean HTTP endpoints with JSON responses
- **Testing**: Comprehensive unit and integration tests

## Tech Stack

- **Language**: Go 1.25.1
- **Web Framework**: Chi router
- **Database**: PostgreSQL (with lib/pq driver)
- **Architecture**: Clean Architecture (Domain-driven design)
- **Testing**: Go testing package
- **Documentation**: Swagger/OpenAPI

## Project Structure

```
src/
├── cmd/api/           # Application entry point
├── internal/
│   ├── api/           # Route registration
│   ├── domain/        # Domain interfaces
│   ├── handler/       # HTTP handlers
│   ├── model/         # Data models
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── test/          # Tests
├── fitness-db/        # Database schema and setup
└── util/              # Utilities and constants
```

## Quick Start

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL (optional - falls back to in-memory)
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd WorkoutPal-Backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL (optional):
```bash
cd src/fitness-db
docker-compose up -d
```

4. Run the application:
```bash
go run src/cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Testing

Run all tests:
```bash
go test ./src/internal/test/... -v
```

Run specific test files:
```bash
go test ./src/internal/test/handler_test.go -v
go test ./src/internal/test/repository_test.go -v
```

## API Endpoints

### Users
- `GET /users` - Get all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create new user
- `PATCH /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Goals
- `POST /users/{id}/goals` - Create user goal
- `GET /users/{id}/goals` - Get user goals

### Social
- `POST /users/{id}/follow` - Follow user
- `GET /users/{id}/followers` - Get user followers
- `GET /users/{id}/following` - Get users being followed

### Routines
- `POST /users/{id}/routines` - Create workout routine
- `GET /users/{id}/routines` - Get user routines

## Database Schema

The application uses PostgreSQL with the following main tables:
- `users` - User profiles and authentication
- `goals` - User fitness goals
- `follows` - User following relationships
- `workout_routine` - Custom workout routines
- `exercises` - Exercise database

See `src/fitness-db/schema.sql` for complete schema.

## Configuration

The application automatically detects PostgreSQL availability:
- **PostgreSQL available**: Uses PostgreSQL database
- **PostgreSQL unavailable**: Falls back to in-memory storage (ideal for testing)

Database connection: `host=localhost port=5432 user=user password=password dbname=workoutpal`

## Documentation

📋 **[Sprint 0 Documentation](./sprint0.md)** - Complete project overview, features, architecture, and planning details

📋 **[API Documentation](./support_files/Endpoint%20Documentation.md)** - Detailed endpoint documentation

## Swagger Documentation

### Viewing API Documentation

1. Start the server:
```bash
go run src/cmd/api/main.go
```

2. Open your browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

### Regenerating Documentation

1. Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate swagger files:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
swag init -g src/cmd/api/main.go -o src/internal/api/docs
```

3. Restart the server to see updated documentation

## Related Repositories

- [Frontend Repository](https://github.com/Onyelechie/WorkoutPal-Frontend) - React.js frontend application

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./src/internal/test/... -v`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.