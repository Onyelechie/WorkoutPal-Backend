package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type relationshipHandler struct{}

func NewMockRelationshipHandler() handler.RelationshipHandler {
	return &relationshipHandler{}
}

func (h *relationshipHandler) ReadFollowers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.User{user})
}

func (h *relationshipHandler) ReadFollowings(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.User{user})
}

func (h *relationshipHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "success"})
}

func (h *relationshipHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "success"})
}
