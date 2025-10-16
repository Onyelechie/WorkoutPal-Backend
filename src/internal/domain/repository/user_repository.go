package repository

import "workoutpal/src/internal/model"

type UserRepository interface {
	ReadUsers() ([]*model.User, error)
	ReadUserByID(id int64) (*model.User, error)
	ReadUserByEmail(email string) (*model.User, error)
	CreateUser(request model.CreateUserRequest) (*model.User, error)
	UpdateUser(request model.UpdateUserRequest) (*model.User, error)
	DeleteUser(request model.DeleteUserRequest) error
}
