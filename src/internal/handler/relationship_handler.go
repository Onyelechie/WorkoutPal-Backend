package handler

import (
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type relationshipHandler struct {
	relationshipService service.RelationshipService
}

func NewRelationshipHandler(rs service.RelationshipService) handler.RelationshipHandler {
	return &relationshipHandler{
		relationshipService: rs,
	}
}

// ReadFollowers godoc
// @Summary List a user's followers
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} int64 "Follower IDs retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/followers [get]
func (h *relationshipHandler) ReadFollowers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	followers, err := h.relationshipService.ReadUserFollowers(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, followers)
}

// ReadFollowings godoc
// @Summary List users that the target user is following
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} int64 "Following IDs retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/following [get]
func (h *relationshipHandler) ReadFollowings(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	following, err := h.relationshipService.ReadUserFollowing(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, following)
}

// FollowUser godoc
// @Summary Follow a user
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID to follow"
// @Param follower_id query int true "Follower user ID"
// @Success 200 {object} model.BasicResponse "Successfully followed user"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Router /users/{id}/follow [post]
func (h *relationshipHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
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

	err = h.relationshipService.FollowUser(followerID, followeeID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Successfully followed user"})
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID to unfollow"
// @Param follower_id query int true "Follower user ID"
// @Success 200 {object} model.BasicResponse "Successfully unfollowed user"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Router /users/{id}/unfollow [post]
func (h *relationshipHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
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

	err = h.relationshipService.UnfollowUser(followerID, followeeID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Successfully unfollowed user"})
}
