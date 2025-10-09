package repository

import "workoutpal/src/internal/model"

type UserRepository interface {
	ReadUsers() ([]model.User, error)
	GetUserByID(id int64) (model.User, error)
	CreateUser(request model.CreateUserRequest) (model.User, error)
	UpdateUser(request model.UpdateUserRequest) (model.User, error)
	DeleteUser(request model.DeleteUserRequest) error
	// Goals
	CreateGoal(userID int64, request model.CreateGoalRequest) (model.Goal, error)
	GetUserGoals(userID int64) ([]model.Goal, error)
	UpdateGoal(request model.UpdateGoalRequest) (model.Goal, error)
	DeleteGoal(goalID int64) error
	// Followers
	FollowUser(followerID, followeeID int64) error
	UnfollowUser(followerID, followeeID int64) error
	GetUserFollowers(userID int64) ([]int64, error)
	GetUserFollowing(userID int64) ([]int64, error)
	// Routines
	CreateRoutine(userID int64, request model.CreateRoutineRequest) (model.ExerciseRoutine, error)
	GetUserRoutines(userID int64) ([]model.ExerciseRoutine, error)
	DeleteRoutine(routineID int64) error
}
