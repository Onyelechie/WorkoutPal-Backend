package repository

import (
	"database/sql"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &PostRepository{db: db}
}

func (p *PostRepository) ReadPostsByUserID(targetUserID int64, userID int64) ([]*model.Post, error) {
	rows, err := p.db.Query(`
    SELECT 
        p.id,
        p.title,
        p.body,
        p.caption,
        p.status,
        p.created_at,
        u.username,
        COUNT(DISTINCT pl_all.user_id) AS likes,
        pl_user.post_id IS NOT NULL AS is_liked
    FROM posts p 
    LEFT JOIN post_likes pl_all ON p.id = pl_all.post_id
    LEFT JOIN post_likes pl_user ON p.id = pl_user.post_id AND pl_user.user_id = $1
    JOIN users u ON u.id = p.user_id
    WHERE p.user_id = $2
    GROUP BY p.id, u.username, pl_user.post_id`,
		userID, targetUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result = make([]*model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.Caption,
			&post.Status,
			&post.Date,
			&post.PostedBy,
			&post.Likes,
			&post.IsLiked,
		); err != nil {
			return nil, err
		}
		result = append(result, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostRepository) ReadPosts(userID int64) ([]*model.Post, error) {
	rows, err := p.db.Query(`
    SELECT 
        p.id,
        p.title,
        p.body,
        p.caption,
        p.status,
        p.created_at,
        u.username,
        COUNT(DISTINCT pl_all.user_id) AS likes,
        pl_user.post_id IS NOT NULL AS is_liked
    FROM posts p 
    LEFT JOIN post_likes pl_all ON p.id = pl_all.post_id
    LEFT JOIN post_likes pl_user ON p.id = pl_user.post_id AND pl_user.user_id = $1
    JOIN users u ON u.id = p.user_id
    GROUP BY p.id, u.username, pl_user.post_id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Post = make([]*model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.Caption,
			&post.Status,
			&post.Date,
			&post.PostedBy,
			&post.Likes,
			&post.IsLiked,
		); err != nil {
			return nil, err
		}
		result = append(result, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostRepository) ReadPost(id int64, userID int64) (*model.Post, error) {
	row := p.db.QueryRow(`
    SELECT 
        p.id,
        p.title,
        p.body,
        p.caption,
        p.status,
        p.created_at,
        u.username,
        COUNT(DISTINCT pl_all.user_id) AS likes,
        pl_user.post_id IS NOT NULL AS is_liked
    FROM posts p 
    LEFT JOIN post_likes pl_all ON p.id = pl_all.post_id
    LEFT JOIN post_likes pl_user 
        ON p.id = pl_user.post_id AND pl_user.user_id = $1
    JOIN users u 
        ON u.id = p.user_id 
    WHERE p.id = $2
    GROUP BY p.id, u.username, pl_user.post_id`,
		userID, id,
	)

	var post model.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.Caption,
		&post.Status,
		&post.Date,
		&post.PostedBy,
		&post.Likes,
		&post.IsLiked,
	)

	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostRepository) CreatePost(req model.CreatePostRequest) (*model.Post, error) {
	var id int64
	row := p.db.QueryRow(`
		INSERT INTO posts (title, body, caption, status, user_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id`,
		req.Title, req.Body, req.Caption, req.Status, req.PostedBy)

	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return p.ReadPost(id, req.PostedBy)
}

func (p *PostRepository) UpdatePost(req model.UpdatePostRequest) (*model.Post, error) {
	_, err := p.db.Exec(`
		UPDATE posts 
		SET title=$1, body=$2, caption=$3, status=$4 
		WHERE id=$5`,
		req.Title, req.Body, req.Caption, req.Status, req.ID)
	if err != nil {
		return nil, err
	}

	return p.ReadPost(req.ID, 0)
}

func (p *PostRepository) DeletePost(id int64) error {
	_, err := p.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	return err
}

func (p *PostRepository) ReadCommentsByPost(postID int64) ([]*model.Comment, error) {
	rows, err := p.db.Query(`
		SELECT pc.id, pc.body, pc.created_at, u.username 
		FROM post_comments pc 
		JOIN users u ON u.id = pc.user_id 
		WHERE pc.post_id = $1 AND pc.parent_comment_id IS NULL
		ORDER BY pc.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result = make([]*model.Comment, 0)
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.Comment, &c.Date, &c.Username)
		if err != nil {
			return nil, err
		}
		result = append(result, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostRepository) ReadCommentsByComment(commentID int64) ([]*model.Comment, error) {
	rows, err := p.db.Query(`
		SELECT pc.id, pc.body, pc.created_at, u.username 
		FROM post_comments pc 
		JOIN users u ON u.id = pc.user_id 
		WHERE pc.parent_comment_id = $1
		ORDER BY pc.created_at ASC`, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.Comment, &c.Date, &c.Username)
		if err != nil {
			return nil, err
		}
		result = append(result, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostRepository) CommentOnPost(req model.CommentOnPostRequest) error {
	_, err := p.db.Exec(`
		INSERT INTO post_comments (body, user_id, post_id)
		VALUES ($1,$2,$3)`,
		req.Comment, req.UserID, req.PostID)
	return err
}

func (p *PostRepository) CommentOnComment(req model.CommentOnCommentRequest) error {
	_, err := p.db.Exec(`
		INSERT INTO post_comments (body, user_id, post_id, parent_comment_id)
		VALUES ($1,$2,$3,$4)`,
		req.Comment, req.UserID, req.PostID, req.CommentID)
	return err
}

func (p *PostRepository) LikePost(req model.LikePostRequest) (*model.Post, error) {
	_, err := p.db.Exec(`INSERT INTO post_likes(user_id, post_id,created_at) VALUES ($1,$2,now())`, req.UserID, req.PostID)
	if err != nil {
		return nil, err
	}

	return p.ReadPost(req.PostID, req.UserID)
}

func (p *PostRepository) UnlikePost(req model.UnikePostRequest) (*model.Post, error) {
	_, err := p.db.Exec(`DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2`, req.UserID, req.PostID)
	if err != nil {
		return nil, err
	}

	return p.ReadPost(req.PostID, req.UserID)
}
