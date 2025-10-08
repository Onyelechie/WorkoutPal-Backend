# User Profiles Backend Implementation

## Implementation Summary

I've implemented a complete user profiles backend with the following components:

### 1. User Service (`src/internal/service/user_service.go`)
- In-memory storage with thread-safe operations using sync.RWMutex
- CRUD operations: Create, Read (all/by ID), Update, Delete
- Validation for duplicate usernames and emails
- Easy to replace with database implementation later

### 2. User Handler (`src/internal/handler/user_handler.go`)
- RESTful API endpoints with proper HTTP status codes
- JSON request/response handling
- Input validation and error handling
- Swagger documentation annotations

### 3. API Endpoints
- `POST /users` - Create new user
- `GET /users` - Get all users
- `GET /users/{id}` - Get user by ID
- `PATCH /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### 4. Features Implemented
- User profile creation with required fields validation
- Profile updates with conflict checking
- User retrieval (individual and all users)
- User deletion
- Thread-safe in-memory storage
- Proper error responses

## Testing the API

### Create User
```bash
curl -X POST http://127.0.0.1:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "name": "John Doe",
    "age": 25,
    "height": 175.5,
    "heightMetric": "cm",
    "weight": 70.0,
    "weightMetric": "kg"
  }'
```

### Get All Users
```bash
curl http://127.0.0.1:8080/users
```

### Get User by ID
```bash
curl http://127.0.0.1:8080/users/1
```

### Update User
```bash
curl -X PATCH http://127.0.0.1:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe_updated",
    "email": "john.updated@example.com",
    "name": "John Doe Updated",
    "age": 26,
    "height": 176.0,
    "heightMetric": "cm",
    "weight": 72.0,
    "weightMetric": "kg"
  }'
```

### Delete User
```bash
curl -X DELETE http://127.0.0.1:8080/users/1
```

## Next Steps

1. **Database Integration**: Replace in-memory storage with MySQL database operations
2. **Password Hashing**: Implement proper password hashing (bcrypt)
3. **Authentication**: Add JWT token validation
4. **Validation**: Add more comprehensive input validation
5. **Testing**: Add unit and integration tests

## Files Modified/Created

- `src/internal/service/user_service.go` (new)
- `src/internal/handler/user_handler.go` (updated)
- `src/internal/domain/handler/user_handler.go` (updated)
- `src/mock_internal/mock_handler/user_handler.go` (updated)
- `src/internal/api/routes.go` (updated)

The implementation follows the existing project structure and can be easily extended with database connectivity when the database team provides the schema and connection details.