package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type relationshipService struct {
	relationshipRepository repository.RelationshipRepository
	userRepository         repository.UserRepository
}

func NewRelationshipService(relationshipRepository repository.RelationshipRepository, userRepository repository.UserRepository) service.RelationshipService {
	return &relationshipService{
		relationshipRepository: relationshipRepository,
		userRepository:         userRepository,
	}
}

func (u *relationshipService) FollowUser(followerID, followeeID int64) error {
	return u.relationshipRepository.FollowUser(followerID, followeeID)
}

func (u *relationshipService) UnfollowUser(followerID, followeeID int64) error {
	return u.relationshipRepository.UnfollowUser(followerID, followeeID)
}

func (u *relationshipService) ReadUserFollowers(userID int64) ([]model.User, error) {
	followerIDs, err := u.relationshipRepository.ReadUserFollowers(userID)
	if err != nil {
		return nil, err
	}

	var followers []model.User
	for _, id := range followerIDs {
		user, err := u.userRepository.ReadUserByID(id)
		if err == nil {
			followers = append(followers, *user)
		}
	}
	return followers, nil
}

func (u *relationshipService) ReadUserFollowing(userID int64) ([]model.User, error) {
	followingIDs, err := u.relationshipRepository.ReadUserFollowing(userID)
	if err != nil {
		return nil, err
	}

	var following []model.User
	for _, id := range followingIDs {
		user, err := u.userRepository.ReadUserByID(id)
		if err == nil {
			following = append(following, *user)
		}
	}
	return following, nil
}
