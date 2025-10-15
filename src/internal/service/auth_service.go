package service

import (
	"context"
	"errors"
	"fmt"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(ur repository.UserRepository) service.AuthService {
	return &authService{userRepository: ur}
}

func (a *authService) Authenticate(ctx context.Context, request model.LoginRequest) (*model.User, error) {

	user, err := a.userRepository.ReadUserByEmail(request.Email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		fmt.Println(user.Password)
		return nil, errors.New("invalid password")
	}

	return user, nil
}
