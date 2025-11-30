package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
				// Return basic user info even for private profiles
				privateProfileResponse := map[string]interface{}{
					"message":  "This profile is private",
					"isPrivate": true,
					"user": map[string]interface{}{
						"id":       user.ID,
						"name":     user.Name,
						"username": user.Username,
						"avatar":   user.Avatar,
					},
				}
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, privateProfileResponse)
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

// UploadAvatar godoc
// @Summary Upload user avatar as base64
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param avatar body map[string]string true "Avatar base64 data"
// @Success 200 {object} map[string]string "Avatar uploaded successfully"
// @Failure 400 {object} model.BasicResponse "Invalid data or user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id}/avatar [post]
func (u *userHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.ID_KEY).(int64)
	viewerID, ok := r.Context().Value(constants.USER_ID_KEY).(int64)
	
	// Debug: Check if viewerID was properly extracted
	if !ok {
		responseErr := util.Error(errors.New("Authentication failed: no user ID in context"), r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Only allow users to upload their own avatar
	if viewerID != userID {
		responseErr := util.Error(errors.New("You can only upload your own avatar"), r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Parse JSON request
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.ErrorResponse(w, r, util.Error(err, r.URL.Path))
		return
	}

	avatarData, ok := req["avatar"]
	if !ok || avatarData == "" {
		responseErr := util.Error(errors.New("Avatar data is required"), r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Validate base64 format (should start with data:image/)
	if !strings.HasPrefix(avatarData, "data:image/") {
		responseErr := util.Error(errors.New("Invalid image format. Must be base64 encoded image"), r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Validate size (approximate - base64 is ~33% larger than original)
	if len(avatarData) > 7000000 { // ~5MB original file
		responseErr := util.Error(errors.New("Image size too large. Must be less than 5MB"), r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Get current user data to preserve other fields
	currentUser, err := u.userService.ReadUserByID(userID)
	if err != nil {
		util.ErrorResponse(w, r, util.Error(err, r.URL.Path))
		return
	}

	// Update user's avatar data in database
	updateReq := model.UpdateUserRequest{
		ID:     userID,
		Avatar: avatarData,
	}

	// Preserve current user data
	updateReq.Username = currentUser.Username
	updateReq.Name = currentUser.Name
	updateReq.Email = currentUser.Email
	updateReq.Age = currentUser.Age
	updateReq.Height = currentUser.Height
	updateReq.HeightMetric = currentUser.HeightMetric
	updateReq.Weight = currentUser.Weight
	updateReq.WeightMetric = currentUser.WeightMetric
	updateReq.IsPrivate = currentUser.IsPrivate
	updateReq.ShowMetricsToFollowers = currentUser.ShowMetricsToFollowers

	updatedUser, err := u.userService.UpdateUser(updateReq)
	if err != nil {
		util.ErrorResponse(w, r, util.Error(err, r.URL.Path))
		return
	}

	render.JSON(w, r, map[string]string{
		"message": "Avatar uploaded successfully",
		"avatar":  updatedUser.Avatar,
	})
}
