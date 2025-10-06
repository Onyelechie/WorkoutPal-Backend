package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type postHandler struct{}

func NewMockPostHandler() handler.PostHandler {
	return &postHandler{}
}

func (p *postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, post)
}

func (p *postHandler) CommentOnPost(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "success"})
}

func (p *postHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "success"})
}

func (p *postHandler) ReadPosts(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.Post{post})
}
