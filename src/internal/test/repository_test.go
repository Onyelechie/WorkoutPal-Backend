package test

import (
	"testing"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
)

func TestUserRepository_CreateUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	req := model.CreateUserRequest{
		Username:     "testuser",
		Email:        "test@example.com",
		Name:         "Test User",
		Password:     "password123",
		Height:       175.5,
		HeightMetric: "cm",
		Weight:       70.0,
		WeightMetric: "kg",
	}

	user, err := repo.CreateUser(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, user.Username)
	}
	if user.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, user.Email)
	}
	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user first
	req := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	createdUser, _ := repo.CreateUser(req)

	// Get user by ID
	user, err := repo.ReadUserByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, user.Username)
	}
}

func TestUserRepository_ReadUsers(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create multiple users
	users := []model.CreateUserRequest{
		{Username: "user1", Email: "user1@example.com", Name: "User 1"},
		{Username: "user2", Email: "user2@example.com", Name: "User 2"},
	}

	for _, req := range users {
		repo.CreateUser(req)
	}

	allUsers, err := repo.ReadUsers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(allUsers) != 2 {
		t.Errorf("Expected 2 users, got %d", len(allUsers))
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user first
	createReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	createdUser, _ := repo.CreateUser(createReq)

	// Update user
	updateReq := model.UpdateUserRequest{
		ID:       createdUser.ID,
		Username: "updateduser",
		Email:    "updated@example.com",
		Name:     "Updated User",
	}

	updatedUser, err := repo.UpdateUser(updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedUser.Username != updateReq.Username {
		t.Errorf("Expected username %s, got %s", updateReq.Username, updatedUser.Username)
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user first
	req := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	createdUser, _ := repo.CreateUser(req)

	// Delete user
	err := repo.DeleteUser(*model.DeleteUserRequest{ID: createdUser.ID})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify user is deleted
	_, err = repo.ReadUserByID(createdUser.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}

func TestUserRepository_CreateGoal(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user first
	userReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	user, _ := repo.CreateUser(userReq)

	// Create goal
	goalReq := model.CreateGoalRequest{
		Name:        "Weight Loss Goal",
		Description: "Lose weight to 65kg",
		Deadline:    "2024-12-31",
	}

	goal, err := repo.CreateGoal(user.ID, goalReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.Name != goalReq.Name {
		t.Errorf("Expected goal name %s, got %s", goalReq.Name, goal.Name)
	}
	if goal.UserID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, goal.UserID)
	}
}

func TestUserRepository_UpdateGoal(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user and goal first
	user, _ := repo.CreateUser(*model.CreateUserRequest{
		Username: "testuser", Email: "test@example.com", Name: "Test User",
	})
	goal, _ := repo.CreateGoal(user.ID, model.CreateGoalRequest{
		Name: "Original Goal", Description: "Original description", Deadline: "2024-12-31",
	})

	// Update goal
	updateReq := model.UpdateGoalRequest{
		ID: goal.ID, Name: "Updated Goal", Description: "Updated description",
		Deadline: "2025-01-31", Status: "completed",
	}

	updatedGoal, err := repo.UpdateGoal(updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedGoal.Name != updateReq.Name {
		t.Errorf("Expected goal name %s, got %s", updateReq.Name, updatedGoal.Name)
	}
}

func TestUserRepository_DeleteGoal(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user and goal first
	user, _ := repo.CreateUser(*model.CreateUserRequest{
		Username: "testuser", Email: "test@example.com", Name: "Test User",
	})
	goal, _ := repo.CreateGoal(user.ID, model.CreateGoalRequest{
		Name: "Test Goal", Description: "Test description", Deadline: "2024-12-31",
	})

	// Delete goal
	err := repo.DeleteGoal(goal.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify goal is deleted
	goals, _ := repo.GetUserGoals(user.ID)
	if len(goals) != 0 {
		t.Error("Expected goal to be deleted")
	}
}

func TestUserRepository_DeleteRoutine(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create user and routine first
	user, _ := repo.CreateUser(*model.CreateUserRequest{
		Username: "testuser", Email: "test@example.com", Name: "Test User",
	})
	routine, _ := repo.CreateRoutine(user.ID, model.CreateRoutineRequest{
		Name: "Test Routine", Description: "Test description",
	})

	// Delete routine
	err := repo.DeleteRoutine(routine.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify routine is deleted
	routines, _ := repo.GetUserRoutines(user.ID)
	if len(routines) != 0 {
		t.Error("Expected routine to be deleted")
	}
}

func TestUserRepository_FollowUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// Create two users
	user1, _ := repo.CreateUser(*model.CreateUserRequest{
		Username: "user1", Email: "user1@example.com", Name: "User 1",
	})
	user2, _ := repo.CreateUser(*model.CreateUserRequest{
		Username: "user2", Email: "user2@example.com", Name: "User 2",
	})

	// User1 follows User2
	err := repo.FollowUser(user1.ID, user2.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check followers
	followers, err := repo.GetUserFollowers(user2.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(followers) != 1 || followers[0] != user1.ID {
		t.Error("Expected user1 to be following user2")
	}

	// Check following
	following, err := repo.GetUserFollowing(user1.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(following) != 1 || following[0] != user2.ID {
		t.Error("Expected user1 to be following user2")
	}
}
