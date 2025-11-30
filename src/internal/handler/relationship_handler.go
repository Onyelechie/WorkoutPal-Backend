package handler

import (
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/util"

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
// @Success 200 {array} model.User "Followers retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/followers [get]
func (h *relationshipHandler) ReadFollowers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	followers, err := h.relationshipService.ReadUserFollowers(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, followers)
}

// ReadFollowings godoc
// @Summary List users that the target user is following
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} model.User "Following users retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/following [get]
func (h *relationshipHandler) ReadFollowings(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	following, err := h.relationshipService.ReadUserFollowing(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
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
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	followerIDStr := r.URL.Query().Get("follower_id")
	followerID, err := strconv.ParseInt(followerIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	err = h.relationshipService.FollowUser(followerID, followeeID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
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
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	followerIDStr := r.URL.Query().Get("follower_id")
	followerID, err := strconv.ParseInt(followerIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	err = h.relationshipService.UnfollowUser(followerID, followeeID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Successfully unfollowed user"})
}

// SendFollowRequest godoc
// @Summary Send a follow request to a private profile
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID to request to follow"
// @Param requester_id query int true "Requester user ID"
// @Success 200 {object} model.BasicResponse "Follow request sent successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Router /users/{id}/follow-request [post]
func (h *relationshipHandler) SendFollowRequest(w http.ResponseWriter, r *http.Request) {
	requestedIDStr := chi.URLParam(r, "id")
	requestedID, err := strconv.ParseInt(requestedIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	requesterIDStr := r.URL.Query().Get("requester_id")
	requesterID, err := strconv.ParseInt(requesterIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	err = h.relationshipService.SendFollowRequest(requesterID, requestedID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Follow request sent successfully"})
}

// GetPendingFollowRequests godoc
// @Summary Get all pending follow requests for the authenticated user
// @Tags Relationships
// @Produce json
// @Success 200 {array} model.FollowRequestWithUser "Pending follow requests retrieved successfully"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Router /follow-requests [get]
func (h *relationshipHandler) GetPendingFollowRequests(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	requests, err := h.relationshipService.GetPendingFollowRequests(userID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// Return empty array if no requests instead of null
	if requests == nil {
		requests = []*model.FollowRequestWithUser{}
	}

	render.JSON(w, r, requests)
}

// RespondToFollowRequest godoc
// @Summary Accept or reject a follow request
// @Tags Relationships
// @Accept json
// @Produce json
// @Param request body model.FollowRequestResponse true "Follow request response"
// @Success 200 {object} model.BasicResponse "Follow request processed successfully"
// @Failure 400 {object} model.BasicResponse "Invalid request"
// @Router /follow-requests/respond [post]
func (h *relationshipHandler) RespondToFollowRequest(w http.ResponseWriter, r *http.Request) {
	var req model.FollowRequestResponse
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	var err error
	if req.Action == "accept" {
		err = h.relationshipService.AcceptFollowRequest(req.RequestID)
	} else if req.Action == "reject" {
		err = h.relationshipService.RejectFollowRequest(req.RequestID)
	} else {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid action. Use 'accept' or 'reject'"})
		return
	}

	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Follow request processed successfully"})
}

// CancelFollowRequest godoc
// @Summary Cancel a pending follow request
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID whose request to cancel"
// @Param requester_id query int true "Requester user ID"
// @Success 200 {object} model.BasicResponse "Follow request cancelled successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Router /users/{id}/follow-request [delete]
func (h *relationshipHandler) CancelFollowRequest(w http.ResponseWriter, r *http.Request) {
	requestedIDStr := chi.URLParam(r, "id")
	requestedID, err := strconv.ParseInt(requestedIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	requesterIDStr := r.URL.Query().Get("requester_id")
	requesterID, err := strconv.ParseInt(requesterIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	err = h.relationshipService.CancelFollowRequest(requesterID, requestedID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Follow request cancelled successfully"})
}

// GetFollowRequestStatus godoc
// @Summary Get the status of a follow request
// @Tags Relationships
// @Produce json
// @Param id path int true "User ID"
// @Param requester_id query int true "Requester user ID"
// @Success 200 {object} model.FollowRequestModel "Follow request status"
// @Failure 404 {object} model.BasicResponse "No follow request found"
// @Router /users/{id}/follow-request/status [get]
func (h *relationshipHandler) GetFollowRequestStatus(w http.ResponseWriter, r *http.Request) {
	requestedIDStr := chi.URLParam(r, "id")
	requestedID, err := strconv.ParseInt(requestedIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	requesterIDStr := r.URL.Query().Get("requester_id")
	requesterID, err := strconv.ParseInt(requesterIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	req, err := h.relationshipService.GetFollowRequest(requesterID, requestedID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	if req == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: "No follow request found"})
		return
	}

	render.JSON(w, r, req)
}
