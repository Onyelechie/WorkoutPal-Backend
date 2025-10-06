package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
)

type relationshipHandler struct{}

func NewRelationshipHandler() handler.RelationshipHandler {
	return &relationshipHandler{}
}

// ReadFollowers godoc
// @Summary List a user's followers
// @Tags Relationships
// @Accept json
// @Produce json
// @Success 200 {array} model.User "Followers retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id}/followers [get]
func (h *relationshipHandler) ReadFollowers(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// ReadFollowings godoc
// @Summary List users that the target user is following
// @Tags Relationships
// @Accept json
// @Produce json
// @Success 200 {array} model.User "Followings retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /users/{id}/followings [get]
func (h *relationshipHandler) ReadFollowings(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// FollowUser godoc
// @Summary Follow a user
// @Tags Relationships
// @Accept json
// @Produce json
// @Param request body model.FollowRequest true "User to follow (userID)"
// @Success 200 {object} model.BasicResponse "Followed successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /relationships/follow [post]
func (h *relationshipHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Tags Relationships
// @Accept json
// @Produce json
// @Success 200 {object} model.BasicResponse "Unfollowed successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /relationships/unfollow [post]
func (h *relationshipHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
