package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	handler_impl "workoutpal/src/internal/handler"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
	service_impl "workoutpal/src/internal/service"
)

// Relationship handler test helpers
func setupRelationshipHandler() (handler.RelationshipHandler, service.UserService) {
	repo := repository.NewInMemoryUserRepository()
	userService := service_impl.NewUserService(repo)
	relationshipHandler := handler_impl.NewRelationshipHandler(userService)
	return relationshipHandler, userService
}

func createTwoTestUsers(userService service.UserService) {
	user1 := model.CreateUserRequest{Username: "user1", Email: "user1@example.com", Name: "User 1"}
	user2 := model.CreateUserRequest{Username: "user2", Email: "user2@example.com", Name: "User 2"}
	userService.CreateUser(user1)
	userService.CreateUser(user2)
}

func createFollowRelationship(userService service.UserService) {
	userService.FollowUser(1, 2) // user1 follows user2
}

func TestRelationshipHandler_FollowUser(t *testing.T) {
	relationshipHandler, userService := setupRelationshipHandler()
	createTwoTestUsers(userService)

	req := createRequestWithContext("POST", "/users/2/follow?follower_id=1", "2", nil)
	w := httptest.NewRecorder()

	relationshipHandler.FollowUser(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.BasicResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, "Successfully followed user", response.Message, "message")
}

func TestRelationshipHandler_ReadFollowers(t *testing.T) {
	relationshipHandler, userService := setupRelationshipHandler()
	createTwoTestUsers(userService)
	createFollowRelationship(userService)

	req := createRequestWithContext("GET", "/users/2/followers", "2", nil)
	w := httptest.NewRecorder()

	relationshipHandler.ReadFollowers(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response []int64
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 1 || response[0] != 1 {
		t.Errorf("Expected follower ID 1, got %v", response)
	}
}

func TestRelationshipHandler_ReadFollowings(t *testing.T) {
	relationshipHandler, userService := setupRelationshipHandler()
	createTwoTestUsers(userService)
	createFollowRelationship(userService)

	req := createRequestWithContext("GET", "/users/1/following", "1", nil)
	w := httptest.NewRecorder()

	relationshipHandler.ReadFollowings(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response []int64
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 1 || response[0] != 2 {
		t.Errorf("Expected following ID 2, got %v", response)
	}
}