package service

import "workoutpal/src/internal/model"

type PostService interface {
	ReadPostsByUserID(targetUserID int64, userID int64) ([]*model.Post, error)
	ReadPosts(userID int64) ([]*model.Post, error)
	CreatePost(req model.CreatePostRequest) (*model.Post, error)
	UpdatePost(req model.UpdatePostRequest) (*model.Post, error)
	DeletePost(id int64) error

	LikePost(req model.LikePostRequest) (*model.Post, error)
	UnlikePost(req model.UnikePostRequest) (*model.Post, error)

	CommentOnPost(req model.CommentOnPostRequest) error
	CommentOnComment(req model.CommentOnCommentRequest) error
}
