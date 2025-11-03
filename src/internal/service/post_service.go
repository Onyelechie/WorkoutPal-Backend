package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) service.PostService {
	return &PostService{repo: repo}
}

func (s *PostService) ReadPostsByUserId(userID int64) ([]*model.Post, error) {
	return s.repo.ReadPostsByUserId(userID)
}

func (s *PostService) CreatePost(req model.CreatePostRequest) (*model.Post, error) {
	post, err := s.repo.CreatePost(req)
	if err != nil {
		return nil, err
	}

	comments, err := s.repo.ReadCommentsByPost(post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments
	return post, nil
}

func (s *PostService) ReadPosts() ([]*model.Post, error) {
	posts, err := s.repo.ReadPosts()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		comments, err := s.repo.ReadCommentsByPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
	}
	return posts, nil
}

func (s *PostService) UpdatePost(req model.UpdatePostRequest) (*model.Post, error) {
	post, err := s.repo.UpdatePost(req)
	if err != nil {
		return nil, err
	}

	comments, err := s.repo.ReadCommentsByPost(post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments
	return post, nil
}

func (s *PostService) DeletePost(id int64) error {
	return s.repo.DeletePost(id)
}

func (s *PostService) CommentOnPost(req model.CommentOnPostRequest) error {
	return s.repo.CommentOnPost(req)
}

func (s *PostService) CommentOnComment(req model.CommentOnCommentRequest) error {
	return s.repo.CommentOnComment(req)
}
