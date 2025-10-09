package service

import "workoutpal/src/internal/model"

type UserService interface {
	ReadUsers() ([]model.User, error)
	GetUserByID(id int64) (model.User, error)
	CreateUser(request model.CreateUserRequest) (model.User, error)
	UpdateUser(request model.UpdateUserRequest) (model.User, error)
	DeleteUser(request model.DeleteUserRequest) error
	// Goals
	CreateGoal(userID int64, request model.CreateGoalRequest) (model.Goal, error)
	GetUserGoals(userID int64) ([]model.Goal, error)
	// Followers
	FollowUser(followerID, followeeID int64) error
	UnfollowUser(followerID, followeeID int64) error
	GetUserFollowers(userID int64) ([]int64, error)
	GetUserFollowing(userID int64) ([]int64, error)
	// Routines
	CreateRoutine(userID int64, request model.CreateRoutineRequest) (model.ExerciseRoutine, error)
	GetUserRoutines(userID int64) ([]model.ExerciseRoutine, error)
}
