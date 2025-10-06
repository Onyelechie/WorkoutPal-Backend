package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
)

type postHandler struct{}

func NewPostHandler() handler.PostHandler {
	return &postHandler{}
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
func (p *postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
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
func (p *postHandler) CommentOnPost(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// LikePost godoc
// @Summary Like a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body model.LikePostRequest true "Target post to like"
// @Success 200 {object} model.BasicResponse "Post liked successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /posts/like [post]
func (p *postHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
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
func (p *postHandler) ReadPosts(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
