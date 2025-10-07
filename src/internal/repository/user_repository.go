package repository

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type userRepository struct {
	// TODO db connection here
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (u *userRepository) ReadUsers() ([]model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) CreateUser() (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) UpdateUser() (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) DeleteUser() (model.User, error) {
	//TODO implement me
	panic("implement me")
}
