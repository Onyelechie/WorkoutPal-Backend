package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type userHandler struct{}

func NewMockUserHandler() handler.UserHandler {
	return &userHandler{}
}

func (u *userHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (u *userHandler) ReadAllUsers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.User{user})
}

func (u *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (u *userHandler) CreateUserGoal(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "Goal created"})
}

func (u *userHandler) GetUserGoals(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.Goal{})
}

func (u *userHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "User followed"})
}

func (u *userHandler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []int64{})
}

func (u *userHandler) GetUserFollowing(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []int64{})
}

func (u *userHandler) CreateUserRoutine(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "Routine created"})
}

func (u *userHandler) GetUserRoutines(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.ExerciseRoutine{})
}

func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}
