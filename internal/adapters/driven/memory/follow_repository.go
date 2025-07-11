package memory

import (
	"sync"
	"tuity/internal/core/domain"
	"tuity/pkg/errors"
)

type FollowMemoryRepository struct {
	follows         map[string]*domain.Follow
	followsByUser   map[string]map[string]*domain.Follow // followerID -> followedID -> Follow
	followersByUser map[string]map[string]*domain.Follow // followedID -> followerID -> Follow
	mutex           sync.RWMutex
}

func NewFollowMemoryRepository() *FollowMemoryRepository {
	return &FollowMemoryRepository{
		follows:         make(map[string]*domain.Follow),
		followsByUser:   make(map[string]map[string]*domain.Follow),
		followersByUser: make(map[string]map[string]*domain.Follow),
	}
}

func (r *FollowMemoryRepository) Save(follow *domain.Follow) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.follows[follow.ID] = follow

	if r.followsByUser[follow.FollowerID] == nil {
		r.followsByUser[follow.FollowerID] = make(map[string]*domain.Follow)
	}
	if r.followersByUser[follow.FollowedID] == nil {
		r.followersByUser[follow.FollowedID] = make(map[string]*domain.Follow)
	}

	r.followsByUser[follow.FollowerID][follow.FollowedID] = follow
	r.followersByUser[follow.FollowedID][follow.FollowerID] = follow

	return nil
}

func (r *FollowMemoryRepository) FindByID(id string) (*domain.Follow, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	follow, exists := r.follows[id]
	if !exists {
		return nil, errors.NewNotFoundError("follow")
	}
	return follow, nil
}

func (r *FollowMemoryRepository) FindByFollowerAndFollowed(followerID, followedID string) (*domain.Follow, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if userFollows, exists := r.followsByUser[followerID]; exists {
		if follow, exists := userFollows[followedID]; exists {
			return follow, nil
		}
	}
	return nil, errors.NewNotFoundError("follow")
}

func (r *FollowMemoryRepository) FindFollowing(followerID string) ([]*domain.Follow, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userFollows, exists := r.followsByUser[followerID]
	if !exists {
		return []*domain.Follow{}, nil
	}

	result := make([]*domain.Follow, 0, len(userFollows))
	for _, follow := range userFollows {
		result = append(result, follow)
	}
	return result, nil
}

func (r *FollowMemoryRepository) FindFollowers(followedID string) ([]*domain.Follow, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userFollowers, exists := r.followersByUser[followedID]
	if !exists {
		return []*domain.Follow{}, nil
	}

	result := make([]*domain.Follow, 0, len(userFollowers))
	for _, follow := range userFollowers {
		result = append(result, follow)
	}
	return result, nil
}

func (r *FollowMemoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	follow, exists := r.follows[id]
	if !exists {
		return errors.NewNotFoundError("follow")
	}

	delete(r.follows, id)

	if userFollows, exists := r.followsByUser[follow.FollowerID]; exists {
		delete(userFollows, follow.FollowedID)
		if len(userFollows) == 0 {
			delete(r.followsByUser, follow.FollowerID)
		}
	}

	if userFollowers, exists := r.followersByUser[follow.FollowedID]; exists {
		delete(userFollowers, follow.FollowerID)
		if len(userFollowers) == 0 {
			delete(r.followersByUser, follow.FollowedID)
		}
	}

	return nil
}
