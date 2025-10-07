package service

import "workoutpal/src/internal/model"

type UserService interface {
	ReadUsers() ([]model.User, error)
	CreateUser(request model.CreateUserRequest) (model.User, error)
	UpdateUser(request model.UpdateUserRequest) (model.User, error)
	DeleteUser(request model.DeleteUserRequest) (model.User, error)
}
