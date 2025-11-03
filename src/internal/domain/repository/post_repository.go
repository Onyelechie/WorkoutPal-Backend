package repository

import "workoutpal/src/internal/model"

type PostRepository interface {
	ReadPostsByUserId(userID int64) ([]*model.Post, error)
	ReadPosts() ([]*model.Post, error)
	ReadPost(id int64) (*model.Post, error)
	CreatePost(req model.CreatePostRequest) (*model.Post, error)
	UpdatePost(req model.UpdatePostRequest) (*model.Post, error)
	DeletePost(id int64) error

	ReadCommentsByPost(id int64) ([]*model.Comment, error)
	ReadCommentsByComment(id int64) ([]*model.Comment, error)
	CommentOnPost(req model.CommentOnPostRequest) error
	CommentOnComment(req model.CommentOnCommentRequest) error
}
