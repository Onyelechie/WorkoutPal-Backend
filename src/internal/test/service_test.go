package test

import (
	"testing"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
	"workoutpal/src/internal/service"
)

func TestUserService_CreateUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	req := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	user, err := userService.CreateUser(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, user.Username)
	}
}

func TestUserService_ReadUsers(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	// Create test users
	users := []model.CreateUserRequest{
		{Username: "user1", Email: "user1@example.com", Name: "User 1"},
		{Username: "user2", Email: "user2@example.com", Name: "User 2"},
	}

	for _, req := range users {
		userService.CreateUser(req)
	}

	allUsers, err := userService.ReadUsers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	req := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	createdUser, _ := userService.CreateUser(req)

	user, err := userService.ReadUserByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, user.Username)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	// Create user
	createReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	createdUser, _ := userService.CreateUser(createReq)

	// Update user
	updateReq := model.UpdateUserRequest{
		ID:       createdUser.ID,
		Username: "updateduser",
		Email:    "updated@example.com",
		Name:     "Updated User",
	}

	updatedUser, err := userService.UpdateUser(updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser.Username != updateReq.Username {
		t.Errorf("Expected username %s, got %s", updateReq.Username, updatedUser.Username)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	// Create user
	req := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	createdUser, _ := userService.CreateUser(req)

	// Delete user
	err := userService.DeleteUser(*model.DeleteUserRequest{ID: createdUser.ID})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user is deleted
	_, err = userService.ReadUserByID(createdUser.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}

func TestUserService_CreateGoal(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	// Create user
	userReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	user, _ := userService.CreateUser(userReq)

	// Create goal
	goalReq := model.CreateGoalRequest{
		Name:        "Weight Loss Goal",
		Description: "Lose weight to 65kg",
		Deadline:    "2024-12-31",
	}

	goal, err := userService.CreateGoal(user.ID, goalReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.Name != goalReq.Name {
		t.Errorf("Expected goal name %s, got %s", goalReq.Name, goal.Name)
	}
}

func TestUserService_FollowUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)

	// Create two users
	user1, _ := userService.CreateUser(*model.CreateUserRequest{
		Username: "user1", Email: "user1@example.com", Name: "User 1",
	})
	user2, _ := userService.CreateUser(*model.CreateUserRequest{
		Username: "user2", Email: "user2@example.com", Name: "User 2",
	})

	// User1 follows User2
	err := userService.FollowUser(user1.ID, user2.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify following relationship
	followers, err := userService.ReadUserFollowers(user2.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(followers) != 1 || followers[0] != user1.ID {
		t.Error("Expected user1 to be following user2")
	}
}
