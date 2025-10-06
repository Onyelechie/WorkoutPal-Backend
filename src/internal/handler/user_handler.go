package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type userHandler struct {
	userService *service.UserService
}

func NewUserHandler() handler.UserHandler {
	return &userHandler{
		userService: service.NewUserService(),
	}
}

// CreateNewUser godoc
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "New user payload"
// @Success 201 {object} model.User "User created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users [post]
func (u *userHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Name == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Username, email, and name are required"})
		return
	}

	user, err := u.userService.CreateUser(&req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// ReadAllUsers godoc
// @Summary List users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} model.User "Users retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func (u *userHandler) ReadAllUsers(w http.ResponseWriter, r *http.Request) {
	users := u.userService.GetAllUsers()
	render.JSON(w, r, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User "User retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [get]
func (u *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	user, err := u.userService.GetUserByID(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.UpdateUserRequest true "Updated user payload"
// @Success 200 {object} model.User "User updated successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [patch]
func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	req.ID = id
	user, err := u.userService.UpdateUser(&req)
	if err != nil {
		if err.Error() == "user not found" {
			render.Status(r, http.StatusNotFound)
		} else {
			render.Status(r, http.StatusBadRequest)
		}
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.BasicResponse "User deleted successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	err = u.userService.DeleteUser(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "User deleted successfully"})
}

// CreateUserGoal godoc
// @Summary Create a goal for user
// @Tags Users
// @Param id path int true "User ID"
// @Param request body model.CreateGoalRequest true "Goal data"
// @Success 201 {object} model.Goal
// @Router /users/{id}/goals [post]
func (u *userHandler) CreateUserGoal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	var req model.CreateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	goal, err := u.userService.CreateGoal(id, &req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, goal)
}

// GetUserGoals godoc
// @Summary Get user goals
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {array} model.Goal
// @Router /users/{id}/goals [get]
func (u *userHandler) GetUserGoals(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	goals, err := u.userService.GetUserGoals(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, goals)
}

// FollowUser godoc
// @Summary Follow a user
// @Tags Users
// @Param id path int true "User ID to follow"
// @Param follower_id query int true "Follower user ID"
// @Success 200 {object} model.BasicResponse
// @Router /users/{id}/follow [post]
func (u *userHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	followeeIDStr := chi.URLParam(r, "id")
	followeeID, err := strconv.ParseInt(followeeIDStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	followerIDStr := r.URL.Query().Get("follower_id")
	followerID, err := strconv.ParseInt(followerIDStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid follower ID"})
		return
	}

	err = u.userService.FollowUser(followerID, followeeID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Successfully followed user"})
}

// GetUserFollowers godoc
// @Summary Get user followers
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {array} int64
// @Router /users/{id}/followers [get]
func (u *userHandler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	followers, err := u.userService.GetUserFollowers(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, followers)
}

// GetUserFollowing godoc
// @Summary Get users that this user follows
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {array} int64
// @Router /users/{id}/following [get]
func (u *userHandler) GetUserFollowing(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	following, err := u.userService.GetUserFollowing(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, following)
}

// CreateUserRoutine godoc
// @Summary Create exercise routine for user
// @Tags Users
// @Param id path int true "User ID"
// @Param request body model.CreateRoutineRequest true "Routine data"
// @Success 201 {object} model.ExerciseRoutine
// @Router /users/{id}/routines [post]
func (u *userHandler) CreateUserRoutine(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	var req model.CreateRoutineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	routine, err := u.userService.CreateRoutine(id, &req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, routine)
}

// GetUserRoutines godoc
// @Summary Get user exercise routines
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {array} model.ExerciseRoutine
// @Router /users/{id}/routines [get]
func (u *userHandler) GetUserRoutines(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	routines, err := u.userService.GetUserRoutines(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, routines)
}
