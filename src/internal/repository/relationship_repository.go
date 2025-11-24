package repository

import (
	"database/sql"
	"errors"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
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

// Follow request methods

func (r *relationshipRepository) CreateFollowRequest(requesterID, requestedID int64) error {
	// Try normal insert first
	_, err := r.db.Exec(
		"INSERT INTO follow_requests (requester_id, requested_id, status, created_at, updated_at) VALUES ($1, $2, 'pending', NOW(), NOW())",
		requesterID, requestedID,
	)
	if err == nil {
		return nil
	}

	// On unique violation, inspect existing request and decide behavior.
	existing, getErr := r.GetFollowRequest(requesterID, requestedID)
	if getErr != nil || existing == nil {
		return err // fallback to original error if we can't load existing
	}

	// Helper: check if follow relationship currently exists
	var followExists bool
	row := r.db.QueryRow("SELECT 1 FROM follows WHERE following_user_id = $1 AND followed_user_id = $2 LIMIT 1", requesterID, requestedID)
	var dummy int
	if row.Scan(&dummy) == nil {
		followExists = true
	}

	switch existing.Status {
	case "rejected":
		// Reopen rejected request
		if err2 := r.UpdateFollowRequestStatus(existing.ID, "pending"); err2 != nil {
			return err2
		}
		return nil
	case "accepted":
		// If user unfollowed later (no follow relationship), allow re-open
		if !followExists {
			if err2 := r.UpdateFollowRequestStatus(existing.ID, "pending"); err2 != nil {
				return err2
			}
		}
		return nil
	case "pending":
		// Treat re-request as a "nudge": update the timestamp
		_, _ = r.db.Exec("UPDATE follow_requests SET updated_at = NOW() WHERE id = $1", existing.ID)
		return nil
	default:
		return nil
	}
}

func (r *relationshipRepository) GetFollowRequest(requesterID, requestedID int64) (*model.FollowRequestModel, error) {
	var req model.FollowRequestModel
	err := r.db.QueryRow(
		"SELECT id, requester_id, requested_id, status, created_at, updated_at FROM follow_requests WHERE requester_id = $1 AND requested_id = $2",
		requesterID, requestedID,
	).Scan(&req.ID, &req.RequesterID, &req.RequestedID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *relationshipRepository) GetFollowRequestByID(requestID int64) (*model.FollowRequestModel, error) {
	var req model.FollowRequestModel
	err := r.db.QueryRow(
		"SELECT id, requester_id, requested_id, status, created_at, updated_at FROM follow_requests WHERE id = $1",
		requestID,
	).Scan(&req.ID, &req.RequesterID, &req.RequestedID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *relationshipRepository) GetPendingFollowRequests(userID int64) ([]*model.FollowRequestWithUser, error) {
	rows, err := r.db.Query(`
		SELECT fr.id, fr.requester_id, fr.requested_id, fr.status, fr.created_at,
		       u.id, u.username, u.name, u.email, u.avatar_url
		FROM follow_requests fr
		JOIN users u ON fr.requester_id = u.id
		WHERE fr.requested_id = $1 AND fr.status = 'pending'
		ORDER BY fr.created_at DESC
	`, userID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*model.FollowRequestWithUser
	for rows.Next() {
		var req model.FollowRequestWithUser
		var user model.User
		err := rows.Scan(
			&req.ID, &req.RequesterID, &req.RequestedID, &req.Status, &req.CreatedAt,
			&user.ID, &user.Username, &user.Name, &user.Email, &user.Avatar,
		)
		if err != nil {
			return nil, err
		}
		req.User = &user
		requests = append(requests, &req)
	}
	return requests, nil
}

func (r *relationshipRepository) UpdateFollowRequestStatus(requestID int64, status string) error {
	_, err := r.db.Exec(
		"UPDATE follow_requests SET status = $1, updated_at = NOW() WHERE id = $2",
		status, requestID,
	)
	return err
}

func (r *relationshipRepository) DeleteFollowRequest(requestID int64) error {
	_, err := r.db.Exec("DELETE FROM follow_requests WHERE id = $1", requestID)
	return err
}
