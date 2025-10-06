package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
)

type userHandler struct{}

func NewUserHandler() handler.UserHandler {
	return &userHandler{}
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
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

// UpdateUser godoc
// @Summary Update an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.UpdateUserRequest true "Updated user payload (must include ID)"
// @Success 200 {object} model.User "User updated successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users [put]
func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.DeleteUserRequest true "DeleteUser request containing the user ID"
// @Success 200 {object} model.BasicResponse "User deleted successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users [delete]
func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
