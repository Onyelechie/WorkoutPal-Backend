package handler

import (
	"encoding/json"
	"net/http"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/util"
	"workoutpal/src/util/constants"
)

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// CreatePost godoc
// @Summary Create a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body model.CreatePostRequest true "New post payload"
// @Success 201 {object} model.Post "Post created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /posts [post]
func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	post, err := p.svc.CreatePost(req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

// ReadPosts godoc
// @Summary List posts
// @Tags Posts
// @Accept json
// @Produce json
// @Param followings query bool false "If true, only posts from followed users are returned"
// @Success 200 {array} model.Post "Posts retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /posts [get]
func (p *PostHandler) ReadPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := p.svc.ReadPosts()
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

// CommentOnPost godoc
// @Summary Comment on a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body model.CommentOnPostRequest true "Target post and comment text"
// @Success 200 {object} model.BasicResponse "Comment added successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /posts/comment [post]
func (p *PostHandler) CommentOnPost(w http.ResponseWriter, r *http.Request) {
	var req model.CommentOnPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	if err := p.svc.CommentOnPost(req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	resp := model.BasicResponse{Message: "Success"}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// CommentOnComment godoc
// @Summary Comment on another comment
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body model.CommentOnCommentRequest true "Target comment and reply text"
// @Success 200 {object} model.BasicResponse "Reply added successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /posts/comment/reply [post]
func (p *PostHandler) CommentOnComment(w http.ResponseWriter, r *http.Request) {
	var req model.CommentOnCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	if err := p.svc.CommentOnComment(req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	resp := model.BasicResponse{Message: "Success"}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// DeletePost godoc
// @Summary Delete a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} model.BasicResponse "Post deleted successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 404 {object} model.BasicResponse "Post not found"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /posts/{id} [delete]
func (p *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	if err := p.svc.DeletePost(id); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	resp := model.BasicResponse{Message: "Success"}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
