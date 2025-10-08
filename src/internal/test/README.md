# Test Suite

This directory contains comprehensive tests for the WorkoutPal backend API.

## Test Structure

- **`repository_test.go`** - Tests for data layer (repository pattern)
- **`service_test.go`** - Tests for business logic layer
- **`handler_test.go`** - Tests for HTTP handlers/controllers
- **`integration_test.go`** - End-to-end API workflow tests

## Running Tests

### Run All Tests
```bash
go test ./test/
```

### Run Specific Test File
```bash
go test ./test/ -run TestUserRepository
go test ./test/ -run TestUserService
go test ./test/ -run TestUserHandler
go test ./test/ -run TestIntegration
```

### Run with Verbose Output
```bash
go test ./test/ -v
```

### Run with Coverage
```bash
go test ./test/ -cover
```

## Test Coverage

The test suite covers:

### Repository Layer
- ✅ User CRUD operations
- ✅ Goal management
- ✅ Follow/unfollow functionality
- ✅ Routine management
- ✅ Error handling

### Service Layer
- ✅ Business logic validation
- ✅ Data transformation
- ✅ Error propagation

### Handler Layer
- ✅ HTTP request/response handling
- ✅ JSON serialization/deserialization
- ✅ Status code validation
- ✅ URL parameter parsing

### Integration Tests
- ✅ Complete user workflow
- ✅ Follow system workflow
- ✅ Routine management workflow
- ✅ End-to-end API testing

## Test Database

Tests use the in-memory repository implementation to ensure:
- Fast execution
- No external dependencies
- Isolated test environment
- Consistent test data

## Best Practices

- Each test is independent and isolated
- Tests use descriptive names
- Setup and teardown are handled properly
- Both success and error cases are tested
- Integration tests cover real-world workflows