package repository

import (
	"database/sql"
	"errors"
	"workoutpal/src/internal/domain/repository"
)

type relationshipRepository struct {
	db *sql.DB
}

func NewRelationshipRepository(db *sql.DB) repository.RelationshipRepository {
	return &relationshipRepository{db: db}
}

func (r *relationshipRepository) FollowUser(followerID, followeeID int64) error {
	_, err := r.db.Exec("INSERT INTO follows (following_user_id, followed_user_id, created_at) VALUES ($1, $2, NOW())", followerID, followeeID)
	return err
}

func (r *relationshipRepository) UnfollowUser(followerID, followeeID int64) error {
	result, err := r.db.Exec("DELETE FROM follows WHERE following_user_id = $1 AND followed_user_id = $2", followerID, followeeID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("follow relationship not found")
	}
	return nil
}

func (r *relationshipRepository) ReadUserFollowers(userID int64) ([]int64, error) {
	rows, err := r.db.Query("SELECT following_user_id FROM follows WHERE followed_user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []int64
	for rows.Next() {
		var followerID int64
		err := rows.Scan(&followerID)
		if err != nil {
			return nil, err
		}
		followers = append(followers, followerID)
	}
	return followers, nil
}

func (r *relationshipRepository) ReadUserFollowing(userID int64) ([]int64, error) {
	rows, err := r.db.Query("SELECT followed_user_id FROM follows WHERE following_user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []int64
	for rows.Next() {
		var followedID int64
		err := rows.Scan(&followedID)
		if err != nil {
			return nil, err
		}
		following = append(following, followedID)
	}
	return following, nil
}
