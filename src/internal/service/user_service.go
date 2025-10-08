package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) service.UserService {
	return &userService{
		userRepository: ur,
	}
}

func (u *userService) ReadUsers() ([]model.User, error) {
	return u.userRepository.ReadUsers()
}

func (u *userService) GetUserByID(id int64) (model.User, error) {
	return u.userRepository.GetUserByID(id)
}

func (u *userService) CreateUser(request model.CreateUserRequest) (model.User, error) {
	return u.userRepository.CreateUser(request)
}

func (u *userService) UpdateUser(request model.UpdateUserRequest) (model.User, error) {
	return u.userRepository.UpdateUser(request)
}

func (u *userService) DeleteUser(request model.DeleteUserRequest) error {
	return u.userRepository.DeleteUser(request)
}

func (u *userService) CreateGoal(userID int64, request model.CreateGoalRequest) (model.Goal, error) {
	return u.userRepository.CreateGoal(userID, request)
}

func (u *userService) GetUserGoals(userID int64) ([]model.Goal, error) {
	return u.userRepository.GetUserGoals(userID)
}

func (u *userService) FollowUser(followerID, followeeID int64) error {
	return u.userRepository.FollowUser(followerID, followeeID)
}

func (u *userService) GetUserFollowers(userID int64) ([]int64, error) {
	return u.userRepository.GetUserFollowers(userID)
}

func (u *userService) GetUserFollowing(userID int64) ([]int64, error) {
	return u.userRepository.GetUserFollowing(userID)
}

func (u *userService) CreateRoutine(userID int64, request model.CreateRoutineRequest) (model.ExerciseRoutine, error) {
	return u.userRepository.CreateRoutine(userID, request)
}

func (u *userService) GetUserRoutines(userID int64) ([]model.ExerciseRoutine, error) {
	return u.userRepository.GetUserRoutines(userID)
}
