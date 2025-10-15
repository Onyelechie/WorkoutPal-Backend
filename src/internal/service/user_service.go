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

func (u *userService) ReadUserByEmail(email string) (*model.User, error) {
	return u.userRepository.ReadUserByEmail(email)
}

func (u *userService) ReadUsers() ([]*model.User, error) {
	return u.userRepository.ReadUsers()
}

func (u *userService) ReadUserByID(id int64) (*model.User, error) {
	return u.userRepository.ReadUserByID(id)
}

func (u *userService) CreateUser(request model.CreateUserRequest) (*model.User, error) {
	return u.userRepository.CreateUser(request)
}

func (u *userService) UpdateUser(request model.UpdateUserRequest) (*model.User, error) {
	return u.userRepository.UpdateUser(request)
}

func (u *userService) DeleteUser(request model.DeleteUserRequest) error {
	return u.userRepository.DeleteUser(request)
}
