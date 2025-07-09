package ports

import "tuity/internal/core/domain"

type FollowRepository interface {
	Save(follow *domain.Follow) error
	FindByID(id string) (*domain.Follow, error)
	FindByFollowerAndFollowed(followerID, followedID string) (*domain.Follow, error)
	FindFollowing(followerID string) ([]*domain.Follow, error)
	FindFollowers(followedID string) ([]*domain.Follow, error)
	Delete(id string) error
}
