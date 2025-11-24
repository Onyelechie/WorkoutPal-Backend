package repository

import (
	"errors"
	"regexp"
	"testing"

	"workoutpal/src/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostRepository_CreatePost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.CreatePostRequest{Title: "T", Body: "B", Caption: "C", Status: "active", PostedBy: 1}
	rows := sqlmock.NewRows([]string{"id", "title", "body", "caption", "status", "created_at"}).
		AddRow(1, "T", "B", "C", "active", "now")

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO posts (title, body, caption, status, user_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, title, body, caption, status, created_at`)).
		WithArgs(req.Title, req.Body, req.Caption, req.Status, req.PostedBy).
		WillReturnRows(rows)

	got, err := repo.CreatePost(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.ID != 1 || got.Title != "T" {
		t.Fatalf("unexpected post: %#v", got)
	}
}

func TestPostRepository_CreatePost_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.CreatePostRequest{Title: "Fail"}
	mock.ExpectQuery("INSERT INTO posts").WillReturnError(errors.New("insert fail"))

	_, err := repo.CreatePost(req)
	if err == nil || err.Error() != "insert fail" {
		t.Fatalf("expected insert fail, got %v", err)
	}
}

func TestPostRepository_ReadPosts_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	userID := int64(42)

	rows := sqlmock.NewRows([]string{
		"id", "title", "body", "caption", "status", "created_at", "username", "is_liked",
	}).
		AddRow(1, "A", "B", "C", "active", "now", "user1", true).
		AddRow(2, "X", "Y", "Z", "inactive", "now", "user2", false)

	mock.ExpectQuery("SELECT p.id").
		WithArgs(userID).
		WillReturnRows(rows)

	got, err := repo.ReadPosts(userID)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(got))
	}
	if !got[0].IsLiked {
		t.Fatalf("expected first post to be liked, got IsLiked=%v", got[0].IsLiked)
	}
	if got[1].IsLiked {
		t.Fatalf("expected second post to be not liked, got IsLiked=%v", got[1].IsLiked)
	}
}

func TestPostRepository_ReadPosts_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	userID := int64(42)

	mock.ExpectQuery("SELECT p.id").
		WithArgs(userID).
		WillReturnError(errors.New("fail"))

	_, err := repo.ReadPosts(userID)
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_ReadPost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	postID := int64(1)
	userID := int64(42)

	rows := sqlmock.NewRows([]string{
		"id", "title", "body", "caption", "status", "created_at", "username", "is_liked",
	}).
		AddRow(postID, "T", "B", "C", "active", "now", "user", true)

	mock.ExpectQuery("SELECT p.id").
		WithArgs(userID, postID).
		WillReturnRows(rows)

	got, err := repo.ReadPost(postID, userID)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.ID != postID {
		t.Fatalf("unexpected ID: got %d, want %d", got.ID, postID)
	}
	if !got.IsLiked {
		t.Fatalf("expected IsLiked=true, got %v", got.IsLiked)
	}
}

func TestPostRepository_ReadPost_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	postID := int64(1)
	userID := int64(42)

	mock.ExpectQuery("SELECT p.id").
		WithArgs(userID, postID).
		WillReturnError(errors.New("no rows"))

	_, err := repo.ReadPost(postID, userID)
	if err == nil || err.Error() != "no rows" {
		t.Fatalf("expected no rows, got %v", err)
	}
}

func TestPostRepository_UpdatePost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.UpdatePostRequest{ID: 1, Title: "T", Body: "B", Caption: "C", Status: "S"}
	rows := sqlmock.NewRows([]string{"id", "title", "body", "caption", "status", "created_at"}).
		AddRow(1, "T", "B", "C", "S", "now")

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE posts 
		SET title=$1, body=$2, caption=$3, status=$4 
		WHERE id=$5
		RETURNING id, title, body, caption, status, created_at`)).
		WithArgs(req.Title, req.Body, req.Caption, req.Status, req.ID).
		WillReturnRows(rows)

	got, err := repo.UpdatePost(req)
	if err != nil || got.ID != 1 {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestPostRepository_UpdatePost_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.UpdatePostRequest{ID: 1}
	mock.ExpectQuery("UPDATE posts").WillReturnError(errors.New("fail"))

	_, err := repo.UpdatePost(req)
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_DeletePost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM posts WHERE id = $1")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeletePost(1); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPostRepository_DeletePost_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	mock.ExpectExec("DELETE FROM posts").
		WithArgs(int64(1)).
		WillReturnError(errors.New("fail"))

	if err := repo.DeletePost(1); err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_ReadCommentsByPost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"id", "body", "created_at", "username"}).
		AddRow(1, "C", "now", "user")
	mock.ExpectQuery("SELECT pc.id").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	got, err := repo.ReadCommentsByPost(1)
	if err != nil || len(got) != 1 {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestPostRepository_ReadCommentsByPost_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	mock.ExpectQuery("SELECT pc.id").
		WillReturnError(errors.New("fail"))

	_, err := repo.ReadCommentsByPost(1)
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_ReadCommentsByComment_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"id", "body", "created_at", "username"}).
		AddRow(1, "C", "now", "user")
	mock.ExpectQuery("SELECT pc.id").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	got, err := repo.ReadCommentsByComment(1)
	if err != nil || len(got) != 1 {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestPostRepository_ReadCommentsByComment_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	mock.ExpectQuery("SELECT pc.id").
		WillReturnError(errors.New("fail"))

	_, err := repo.ReadCommentsByComment(1)
	if err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_CommentOnPost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.CommentOnPostRequest{PostID: 1, UserID: 2, Comment: "Nice"}
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO post_comments (body, user_id, post_id)
		VALUES ($1,$2,$3)`)).
		WithArgs(req.Comment, req.UserID, req.PostID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.CommentOnPost(req); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPostRepository_CommentOnPost_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.CommentOnPostRequest{PostID: 1, UserID: 2, Comment: "Bad"}
	mock.ExpectExec("INSERT INTO post_comments").
		WillReturnError(errors.New("fail"))

	if err := repo.CommentOnPost(req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_CommentOnComment_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.CommentOnCommentRequest{CommentID: 1, PostID: 2, UserID: 3, Comment: "Reply"}
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO post_comments (body, user_id, post_id, parent_comment_id)
		VALUES ($1,$2,$3,$4)`)).
		WithArgs(req.Comment, req.UserID, req.PostID, req.CommentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.CommentOnComment(req); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPostRepository_CommentOnComment_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.CommentOnCommentRequest{Comment: "Reply"}
	mock.ExpectExec("INSERT INTO post_comments").
		WillReturnError(errors.New("fail"))

	if err := repo.CommentOnComment(req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostRepository_LikePost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.LikePostRequest{UserID: 2, PostID: 1}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO post_likes(user_id, post_id,created_at) VALUES ($1,$2,now())`)).
		WithArgs(req.UserID, req.PostID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows := sqlmock.NewRows([]string{
		"id", "title", "body", "caption", "status", "created_at", "username", "is_liked",
	}).
		AddRow(1, "T", "B", "C", "active", "now", "user", true)

	mock.ExpectQuery("SELECT p.id").
		WithArgs(req.UserID, req.PostID).
		WillReturnRows(rows)

	got, err := repo.LikePost(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.ID != req.PostID {
		t.Fatalf("expected post ID %d, got %d", req.PostID, got.ID)
	}
	if !got.IsLiked {
		t.Fatalf("expected IsLiked=true after LikePost, got %v", got.IsLiked)
	}
}

func TestPostRepository_LikePost_ErrorOnInsert(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.LikePostRequest{UserID: 2, PostID: 1}

	mock.ExpectExec("INSERT INTO post_likes").
		WithArgs(req.UserID, req.PostID).
		WillReturnError(errors.New("like fail"))

	_, err := repo.LikePost(req)
	if err == nil || err.Error() != "like fail" {
		t.Fatalf("expected like fail, got %v", err)
	}
}

func TestPostRepository_UnlikePost_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.UnikePostRequest{UserID: 2, PostID: 1}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2`)).
		WithArgs(req.UserID, req.PostID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows := sqlmock.NewRows([]string{
		"id", "title", "body", "caption", "status", "created_at", "username", "is_liked",
	}).
		AddRow(1, "T", "B", "C", "active", "now", "user", false)

	mock.ExpectQuery("SELECT p.id").
		WithArgs(req.UserID, req.PostID).
		WillReturnRows(rows)

	got, err := repo.UnlikePost(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.ID != req.PostID {
		t.Fatalf("expected post ID %d, got %d", req.PostID, got.ID)
	}
	if got.IsLiked {
		t.Fatalf("expected IsLiked=false after UnlikePost, got %v", got.IsLiked)
	}
}

func TestPostRepository_UnlikePost_ErrorOnDelete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPostRepository(db)

	req := model.UnikePostRequest{UserID: 2, PostID: 1}

	mock.ExpectExec("DELETE FROM post_likes").
		WithArgs(req.UserID, req.PostID).
		WillReturnError(errors.New("unlike fail"))

	_, err := repo.UnlikePost(req)
	if err == nil || err.Error() != "unlike fail" {
		t.Fatalf("expected unlike fail, got %v", err)
	}
}
