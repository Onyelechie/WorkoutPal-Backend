package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) handler.UserHandler {
	return &userHandler{
		userService: us,
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

	// Validate input
	if err := util.ValidateUsername(req.Username); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}
	if err := util.ValidateEmail(req.Email); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}
	if err := util.ValidateName(req.Name); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}
	if err := util.ValidatePassword(req.Password); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	user, err := u.userService.CreateUser(req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// ReadAllUsers godoc
// @Summary Get all users
// @Tags Users
// @Produce json
// @Success 200 {array} model.User "Users retrieved successfully"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Router /users [get]
func (u *userHandler) ReadAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.ReadUsers()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}
	render.JSON(w, r, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User "User retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
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
// @Summary Update user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.UpdateUserRequest true "Update user payload"
// @Success 200 {object} model.User "User updated successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 404 {object} model.BasicResponse "User not found"
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
	user, err := u.userService.UpdateUser(req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, user)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.BasicResponse "User deleted successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id} [delete]
func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	err = u.userService.DeleteUser(model.DeleteUserRequest{ID: id})
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "User deleted successfully"})
}


