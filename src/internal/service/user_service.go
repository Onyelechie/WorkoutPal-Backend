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
	//TODO implement me
	return u.userRepository.ReadUsers()
}

func (u *userService) CreateUser(request model.CreateUserRequest) (model.User, error) {
	//TODO implement me
	return u.userRepository.CreateUser()
}

func (u *userService) UpdateUser(request model.UpdateUserRequest) (model.User, error) {
	//TODO implement me
	return u.userRepository.UpdateUser()
}

func (u *userService) DeleteUser(request model.DeleteUserRequest) (model.User, error) {
	//TODO implement me
	return u.userRepository.DeleteUser()
}
