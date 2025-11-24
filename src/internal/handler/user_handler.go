package handler

import (
	"encoding/json"
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/util"
	"workoutpal/src/util/constants"

	"github.com/go-chi/render"
)

type userHandler struct {
	userService         service.UserService
	relationshipService service.RelationshipService
}

func NewUserHandler(us service.UserService, rs service.RelationshipService) handler.UserHandler {
	return &userHandler{
		userService:         us,
		relationshipService: rs,
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
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Validate input
	if err := util.ValidateUsername(req.Username); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	if err := util.ValidateEmail(req.Email); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	if err := util.ValidateName(req.Name); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	if err := util.ValidatePassword(req.Password); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	user, err := u.userService.CreateUser(req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
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
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	render.JSON(w, r, users)
}

// ReadUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User "User retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id} [get]
func (u *userHandler) ReadUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)
	viewerID, _ := r.Context().Value(constants.USER_ID_KEY).(int64)
	user, err := u.userService.ReadUserByID(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	// Privacy enforcement: only followers or owner can view private profiles
	if user.IsPrivate && viewerID != 0 && viewerID != id {
		// check if viewer is a follower
		followers, ferr := u.relationshipService.ReadUserFollowers(id)
		if ferr == nil {
			isFollower := false
			for _, f := range followers {
				if f.ID == viewerID {
					isFollower = true
					break
				}
			}
			if !isFollower {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, model.BasicResponse{Message: "This profile is private"})
				return
			}
		}
	}
	// Metrics visibility: only owner or followers (when enabled) can see metrics
	if viewerID != id {
		canSeeMetrics := false
		if user.ShowMetricsToFollowers {
			// Need to be follower
			followers, ferr := u.relationshipService.ReadUserFollowers(id)
			if ferr == nil {
				for _, f := range followers {
					if f.ID == viewerID {
						canSeeMetrics = true
						break
					}
				}
			}
		}
		if !canSeeMetrics {
			// mask metrics for non-authorized viewers
			user.Age = 0
			user.Height = 0
			user.Weight = 0
		}
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
	id := r.Context().Value(constants.ID_KEY).(int64)

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	req.ID = id
	user, err := u.userService.UpdateUser(req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
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
	id := r.Context().Value(constants.ID_KEY).(int64)

	err := u.userService.DeleteUser(model.DeleteUserRequest{ID: id})
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "User deleted successfully"})
}
