package repository

import (
	"database/sql"
	"errors"
	"workoutpal/src/internal/config"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() repository.UserRepository {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		// Fallback to in-memory for testing
		return NewInMemoryUserRepository()
	}
	if err = db.Ping(); err != nil {
		// Fallback to in-memory if DB not available
		return NewInMemoryUserRepository()
	}
	return &userRepository{db: db}
}

// NewPostgresUserRepository forces PostgreSQL usage (no fallback)
func NewPostgresUserRepository() repository.UserRepository {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		panic("Failed to connect to PostgreSQL: " + err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("Failed to ping PostgreSQL: " + err.Error())
	}
	return &userRepository{db: db}
}

func (u *userRepository) ReadUserByEmail(email string) (model.User, error) {
	var user model.User
	var avatarURL sql.NullString
	err := u.db.QueryRow("SELECT id, username, email, password, name, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}
	if avatarURL.Valid {
		user.Avatar = avatarURL.String
	}
	return user, nil
}

func (u *userRepository) ReadUsers() ([]model.User, error) {
	rows, err := u.db.Query("SELECT id, username, email, name, height, height_metric, weight, weight_metric, avatar_url FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var avatarURL sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
		if err != nil {
			return nil, err
		}
		if avatarURL.Valid {
			user.Avatar = avatarURL.String
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *userRepository) GetUserByID(id int64) (model.User, error) {
	var user model.User
	var avatarURL sql.NullString
	err := u.db.QueryRow("SELECT id, username, email, name, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}
	if avatarURL.Valid {
		user.Avatar = avatarURL.String
	}
	return user, nil
}

func (u *userRepository) CreateUser(request model.CreateUserRequest) (model.User, error) {
	var user model.User
	err := u.db.QueryRow(`
		INSERT INTO users (username, email, password, name, height, height_metric, weight, weight_metric, avatar_url) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, username, email, name, height, height_metric, weight, weight_metric, avatar_url`,
		request.Username, request.Email, request.Password, request.Name,
		request.Height, request.HeightMetric, request.Weight, request.WeightMetric, request.Avatar).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name,
		&user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &user.Avatar)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return model.User{}, errors.New("user already exists")
			}
		}
		return model.User{}, err
	}
	return user, nil
}

func (u *userRepository) UpdateUser(request model.UpdateUserRequest) (model.User, error) {
	var user model.User
	var avatarURL sql.NullString
	err := u.db.QueryRow(`
		UPDATE users SET username=$2, email=$3, name=$4, height=$5, height_metric=$6, weight=$7, weight_metric=$8, avatar_url=$9
		WHERE id=$1 RETURNING id, username, email, name, height, height_metric, weight, weight_metric, avatar_url`,
		request.ID, request.Username, request.Email, request.Name,
		request.Height, request.HeightMetric, request.Weight, request.WeightMetric, request.Avatar).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name,
		&user.Height, &user.HeightMetric, &user.Weight, &user.WeightMetric, &avatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}
	if avatarURL.Valid {
		user.Avatar = avatarURL.String
	}
	return user, nil
}

func (u *userRepository) DeleteUser(request model.DeleteUserRequest) error {
	result, err := u.db.Exec("DELETE FROM users WHERE id = $1", request.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (u *userRepository) CreateGoal(userID int64, request model.CreateGoalRequest) (model.Goal, error) {
	var goal model.Goal
	err := u.db.QueryRow(`
		INSERT INTO goals (user_id, name, description, deadline, status) 
		VALUES ($1, $2, $3, $4, 'active') 
		RETURNING id, user_id, name, description, deadline, created_at, status`,
		userID, request.Name, request.Description, request.Deadline).Scan(
		&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.Status)

	if err != nil {
		return model.Goal{}, err
	}
	return goal, nil
}

func (u *userRepository) GetUserGoals(userID int64) ([]model.Goal, error) {
	rows, err := u.db.Query("SELECT id, user_id, name, description, deadline, created_at, status FROM goals WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []model.Goal
	for rows.Next() {
		var goal model.Goal
		err := rows.Scan(&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.Status)
		if err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func (u *userRepository) UpdateGoal(request model.UpdateGoalRequest) (model.Goal, error) {
	var goal model.Goal
	err := u.db.QueryRow(`
		UPDATE goals SET name=$2, description=$3, deadline=$4, status=$5
		WHERE id=$1 RETURNING id, user_id, name, description, deadline, created_at, status`,
		request.ID, request.Name, request.Description, request.Deadline, request.Status).Scan(
		&goal.ID, &goal.UserID, &goal.Name, &goal.Description, &goal.Deadline, &goal.CreatedAt, &goal.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Goal{}, errors.New("goal not found")
		}
		return model.Goal{}, err
	}
	return goal, nil
}

func (u *userRepository) DeleteGoal(goalID int64) error {
	result, err := u.db.Exec("DELETE FROM goals WHERE id = $1", goalID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("goal not found")
	}
	return nil
}

func (u *userRepository) FollowUser(followerID, followeeID int64) error {
	_, err := u.db.Exec("INSERT INTO follows (following_user_id, followed_user_id, created_at) VALUES ($1, $2, NOW())", followerID, followeeID)
	return err
}

func (u *userRepository) UnfollowUser(followerID, followeeID int64) error {
	result, err := u.db.Exec("DELETE FROM follows WHERE following_user_id = $1 AND followed_user_id = $2", followerID, followeeID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("follow relationship not found")
	}
	return nil
}

func (u *userRepository) GetUserFollowers(userID int64) ([]int64, error) {
	rows, err := u.db.Query("SELECT following_user_id FROM follows WHERE followed_user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []int64
	for rows.Next() {
		var followerID int64
		err := rows.Scan(&followerID)
		if err != nil {
			return nil, err
		}
		followers = append(followers, followerID)
	}
	return followers, nil
}

func (u *userRepository) GetUserFollowing(userID int64) ([]int64, error) {
	rows, err := u.db.Query("SELECT followed_user_id FROM follows WHERE following_user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []int64
	for rows.Next() {
		var followedID int64
		err := rows.Scan(&followedID)
		if err != nil {
			return nil, err
		}
		following = append(following, followedID)
	}
	return following, nil
}

func (u *userRepository) CreateRoutine(userID int64, request model.CreateRoutineRequest) (model.ExerciseRoutine, error) {
	var routine model.ExerciseRoutine
	err := u.db.QueryRow(`
		INSERT INTO workout_routine (name, user_id, frequency) 
		VALUES ($1, $2, 'weekly') 
		RETURNING id, name, user_id, frequency`,
		request.Name, userID).Scan(
		&routine.ID, &routine.Name, &routine.UserID, &routine.Description)

	if err != nil {
		return model.ExerciseRoutine{}, err
	}

	// Insert exercises into exercises_in_routine table
	for _, exerciseID := range request.ExerciseIDs {
		_, err := u.db.Exec("INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)", routine.ID, exerciseID)
		if err != nil {
			return model.ExerciseRoutine{}, err
		}
	}

	routine.Description = request.Description
	routine.IsActive = true
	return routine, nil
}

func (u *userRepository) GetUserRoutines(userID int64) ([]model.ExerciseRoutine, error) {
	rows, err := u.db.Query("SELECT id, name, user_id, frequency FROM workout_routine WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routines []model.ExerciseRoutine
	for rows.Next() {
		var routine model.ExerciseRoutine
		err := rows.Scan(&routine.ID, &routine.Name, &routine.UserID, &routine.Description)
		if err != nil {
			return nil, err
		}

		// Fetch exercises for the routine
		exerciseRows, err := u.db.Query("SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1", routine.ID)
		if err != nil {
			return nil, err
		}
		defer exerciseRows.Close()

		var exerciseIDs []int64
		for exerciseRows.Next() {
			var exerciseID int64
			if err := exerciseRows.Scan(&exerciseID); err != nil {
				return nil, err
			}
			exerciseIDs = append(exerciseIDs, exerciseID)
		}
		routine.ExerciseIDs = exerciseIDs

		routine.IsActive = true
		routines = append(routines, routine)
	}
	return routines, nil
}

func (u *userRepository) DeleteRoutine(routineID int64) error {
	result, err := u.db.Exec("DELETE FROM workout_routine WHERE id = $1", routineID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("routine not found")
	}
	return nil
}

func (u *userRepository) GetRoutineWithExercises(routineID int64) (model.ExerciseRoutine, error) {
	var routine model.ExerciseRoutine
	err := u.db.QueryRow("SELECT id, name, user_id, frequency FROM workout_routine WHERE id = $1", routineID).Scan(
		&routine.ID, &routine.Name, &routine.UserID, &routine.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.ExerciseRoutine{}, errors.New("routine not found")
		}
		return model.ExerciseRoutine{}, err
	}
	routine.IsActive = true
	return routine, nil
}

func (u *userRepository) AddExerciseToRoutine(routineID, exerciseID int64) error {
	_, err := u.db.Exec("INSERT INTO routine_exercises (routine_id, exercise_id) VALUES ($1, $2)", routineID, exerciseID)
	return err
}

func (u *userRepository) RemoveExerciseFromRoutine(routineID, exerciseID int64) error {
	result, err := u.db.Exec("DELETE FROM routine_exercises WHERE routine_id = $1 AND exercise_id = $2", routineID, exerciseID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("exercise not found in routine")
	}
	return nil
}
