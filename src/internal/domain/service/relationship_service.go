package service

import "workoutpal/src/internal/model"

type RelationshipService interface {
	FollowUser(followerID, followeeID int64) error
	UnfollowUser(followerID, followeeID int64) error
	ReadUserFollowers(userID int64) ([]model.User, error)
	ReadUserFollowing(userID int64) ([]model.User, error)
}
