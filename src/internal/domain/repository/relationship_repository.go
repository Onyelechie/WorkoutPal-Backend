package repository

type RelationshipRepository interface {
	FollowUser(followerID, followeeID int64) error
	UnfollowUser(followerID, followeeID int64) error
	ReadUserFollowers(userID int64) ([]int64, error)
	ReadUserFollowing(userID int64) ([]int64, error)
}
