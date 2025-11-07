### Users
- `GET /users` - Get all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create new user
- `PATCH /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Goals
- `POST /users/{id}/goals` - Create user goal
- `GET /users/{id}/goals` - Get user goals

### Authentication
- `POST /login` - User login
- `POST /logout` - User logout  
- `GET /me` - Get current authenticated user

### Social/Relationships
- `POST /users/{id}/follow` - Follow user
- `POST /users/{id}/unfollow` - Unfollow user
- `GET /users/{id}/followers` - Get user followers
- `GET /users/{id}/following` - Get users being followed

### Posts
- `GET /posts` - List posts
- `POST /posts` - Create a new post
- `POST /posts/comment` - Comment on a post
- `POST /posts/comment/reply` - Comment on another comment
- `DELETE /posts/{id}` - Delete a post

### User Routines
- `POST /users/{id}/routines` - Create workout routine
- `GET /users/{id}/routines` - Get user routines
- `DELETE /users/{id}/routines/{routine_id}` - Delete user's routine

### Exercises
- `GET /exercises` - Get all exercises

### Routines (Direct Access)
- `GET /routines/{id}` - Get routine with exercises
- `DELETE /routines/{id}` - Delete routine
- `POST /routines/{id}/exercises?exercise_id={exercise_id}` - Add exercise to routine
- `DELETE /routines/{id}/exercises/{exercise_id}` - Remove exercise from routine

### Schedules
`GET /schedules` - Read all schedules for the authenticated user
`GET /schedules/{dayOfWeek}` - Read schedules for the authenticated user on a specific day
`GET /schedules/{id}` - Read a schedule by ID
`POST /schedules` - Create a schedule
`PUT /schedules/{id}` - Update a schedule
`DELETE /schedules/{id}` - Delete a schedule

### Authentication
- `POST /auth/google` - Google OAuth authentication
