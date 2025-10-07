package repository

import "workoutpal/src/internal/model"

// TODO parameter structs need to be defined
type UserRepository interface {
	ReadUsers() ([]model.User, error)
	CreateUser() (model.User, error)
	UpdateUser() (model.User, error)
	DeleteUser() (model.User, error)
}
