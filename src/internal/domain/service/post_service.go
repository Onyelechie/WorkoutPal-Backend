package service

import "workoutpal/src/internal/model"

type PostService interface {
	ReadPosts() ([]*model.Post, error)
	CreatePost(req model.CreatePostRequest) (*model.Post, error)
	UpdatePost(req model.UpdatePostRequest) (*model.Post, error)
	DeletePost(id int64) error

	CommentOnPost(req model.CommentOnPostRequest) error
	CommentOnComment(req model.CommentOnCommentRequest) error
}
