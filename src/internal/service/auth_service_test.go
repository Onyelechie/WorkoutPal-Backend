package service

import (
	"context"
	"errors"
	"testing"
	"workoutpal/src/internal/model"

	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Authenticate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	auth := NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{ID: 1, Email: "test@example.com", Password: string(hashedPassword)}

	mockRepo.
		EXPECT().
		ReadUserByEmail("test@example.com").
		Return(user, nil)

	req := model.LoginRequest{Email: "test@example.com", Password: "password123"}
	result, err := auth.Authenticate(context.Background(), req)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if result == nil || result.Email != user.Email {
		t.Fatalf("unexpected user returned: %+v", result)
	}
}

func TestAuthService_Authenticate_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	auth := NewAuthService(mockRepo)

	mockRepo.
		EXPECT().
		ReadUserByEmail("notfound@example.com").
		Return(nil, errors.New("user not found"))

	req := model.LoginRequest{Email: "notfound@example.com", Password: "irrelevant"}
	_, err := auth.Authenticate(context.Background(), req)

	if err == nil || err.Error() != "invalid email" {
		t.Fatalf("expected 'invalid email', got %v", err)
	}
}

func TestAuthService_Authenticate_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	auth := NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &model.User{ID: 2, Email: "john@example.com", Password: string(hashedPassword)}

	mockRepo.
		EXPECT().
		ReadUserByEmail("john@example.com").
		Return(user, nil)

	req := model.LoginRequest{Email: "john@example.com", Password: "wrongpass"}
	_, err := auth.Authenticate(context.Background(), req)

	if err == nil || err.Error() != "invalid password" {
		t.Fatalf("expected 'invalid password', got %v", err)
	}
}
