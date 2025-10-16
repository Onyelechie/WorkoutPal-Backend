package service

import (
	"context"
	"workoutpal/src/internal/model"
)

type AuthService interface {
	Authenticate(ctx context.Context, request model.LoginRequest) (*model.User, error)
}
