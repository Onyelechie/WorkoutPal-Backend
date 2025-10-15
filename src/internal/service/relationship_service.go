package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
)

type relationshipService struct {
	relationshipRepository repository.RelationshipRepository
}

func NewRelationshipService(relationshipRepository repository.RelationshipRepository) service.RelationshipService {
	return &relationshipService{relationshipRepository: relationshipRepository}
}

func (u *relationshipService) FollowUser(followerID, followeeID int64) error {
	return u.relationshipRepository.FollowUser(followerID, followeeID)
}

func (u *relationshipService) UnfollowUser(followerID, followeeID int64) error {
	return u.relationshipRepository.UnfollowUser(followerID, followeeID)
}

func (u *relationshipService) ReadUserFollowers(userID int64) ([]int64, error) {
	return u.relationshipRepository.ReadUserFollowers(userID)
}

func (u *relationshipService) ReadUserFollowing(userID int64) ([]int64, error) {
	return u.relationshipRepository.ReadUserFollowing(userID)
}
