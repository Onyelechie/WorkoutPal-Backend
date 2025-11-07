**Disclaimer: Lots of ChatGPT code was used in the creation of this project. Particularly for writing tests and repository functions**

# WorkoutPal-Backend

Backend API for the WorkoutPal gym workout tracking application built with Go.

## Project Overview

WorkoutPal is a comprehensive fitness tracking platform that helps users plan, log, and visualize their workout progress. This repository contains the backend API that powers the application.

## Features

- **User Management**: Create, read, update, delete user profiles with age, height, weight tracking
- **Goal Tracking**: Set and monitor fitness goals with deadlines
- **Social Infrastructure**: Backend support for user relationships (followers/following)
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
│   └── service/       # Business logic
├── fitness-db/        # Database schema and setup
├── test/          # Tests
└── util/              # Utilities and constants
```

---

## Quick Start

### Prerequisites

- Go 1.25.1 or higher
- Docker (with docker compose)
- Git

### Installation

1. Download and Install [Go](https://go.dev)

2. Download and Install [Docker](https://www.docker.com/)

3. Clone the repository:
```bash
git clone https://github.com/Onyelechie/WorkoutPal-Backend
cd WorkoutPal-Backend
```

4. Install dependencies:
```bash
go mod download
```

5. Set up PostgreSQL container:
```bash
cd src/fitness-db
docker-compose up -d
```

6. Run the application:
```bash
go run src/cmd/api/main.go
```
**OR Alternatively:**
```bash
# If using Windows Command Prompt
run_go.bat 

# If Windows Powershell
.\run_go.bat 

# If Linux, WSL, MacOS terminal
./run_go.sh 
```

The API will be available at `http://localhost:8080`

---

## Continuous Deployment (CD)

WorkoutPal-Backend uses GitHub Actions to automatically build, push, and deploy Docker images to Azure App Service whenever changes are merged into the main branch.

### How It Works

- Build: On every push to main, GitHub Actions checks out the code and builds a Docker image of the backend.

- Push: The image is pushed to Docker Hub (ilightlysaltedi/workoutpal-backend:latest).

- Deploy: Azure App Service is configured to pull the latest image automatically, updating the running backend without downtime. deployed backend can be found [here](workoutpal-api-daghb9augub5g9ez.canadacentral-01.azurewebsites.net)
---

## Testing

Run all tests with coverage:
```bash
go test ./... -coverpkg=./... -covermode=atomic -coverprofile=coverage.out
```

Run specific test files:
```bash
go test ./src/internal/test/handler_test.go -v
go test ./src/internal/test/repository_test.go -v
```

View coverage in terminal
```bash
go tool cover -func=coverage.out
```

View coverage in browser
```bash
go tool cover -html=coverage.out
```

---

## API Endpoints
[API Endpoints](./support_files/API-Endpoints.md)

## Database Schema

The application uses PostgreSQL with the following main tables:
- `users` - User profiles and authentication
- `goals` - User fitness goals
- `follows` - User following relationships
- `workout_routine` - Custom workout routines
- `exercises` - Exercise database

See `src/fitness-db/schema.sql` for complete schema.

---

## Configuration

The application automatically detects PostgreSQL availability:
- **PostgreSQL available**: Uses PostgreSQL database
- **PostgreSQL unavailable**: Falls back to in-memory storage (ideal for testing)

Database connection: `host=localhost port=5432 user=user password=password dbname=workoutpal`

## Documentation

📋 **[Sprint 0 Documentation](./docs/sprint0.md)** - Complete project overview, features, architecture, and planning details

📋 **[Sprint 1 Documentation](./docs/sprint1.md)** - Testing coverage, test importance, and environment reproducibility.

📋 **[Sprint 2 Documentation](./docs/sprint2.md)** - Regression testing, profiling, test plan and Group 4's best work

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

---

## Regenerating Documentation

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

---

## Manual Testing

### Start the Server
```bash
go run src/cmd/api/main.go
```

---

## Test Scripts
Create and run these curl commands to test the API:

#### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","name":"Test User","password":"password123"}'
```

#### Get All Users
```bash
curl http://localhost:8080/users
```

#### Get All Exercises
```bash
curl http://localhost:8080/exercises
```

#### Create a Goal
```bash
curl -X POST http://localhost:8080/users/1/goals \
  -H "Content-Type: application/json" \
  -d '{"name":"Weight Loss","description":"Lose 10kg","deadline":"2024-12-31"}'
```

#### Create a Routine
```bash
curl -X POST http://localhost:8080/users/1/routines \
  -H "Content-Type: application/json" \
  -d '{"name":"Morning Workout","description":"Daily routine"}'
```

#### Add Exercise to Routine
```bash
curl -X POST "http://localhost:8080/routines/1/exercises?exercise_id=1"
```

#### Follow a User
```bash
curl -X POST "http://localhost:8080/users/2/follow?follower_id=1"
```

#### Get User's Routines
```bash
curl http://localhost:8080/users/1/routines
```

#### Delete User's Routine
```bash
curl -X DELETE http://localhost:8080/users/1/routines/1
```

### Automated Testing Script
Run the complete API test suite:
```bash
./test-api.sh
```

### Testing with Verbose Output
Add `-v` flag to any curl command to see detailed request/response:
```bash
curl -v http://localhost:8080/users
```

### Prerequisites for Testing
- `jq` for JSON formatting: `brew install jq` (macOS) or `apt install jq` (Linux)
- Server running on port 8080

## Related Repositories

- [Frontend Repository](https://github.com/Onyelechie/WorkoutPal-Frontend) - React.js frontend application

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./src/internal/test/... -v`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.