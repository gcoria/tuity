package services

import (
	"tuity/internal/core/domain"
	"tuity/internal/core/ports"
	"tuity/pkg/errors"
	"tuity/pkg/utils"
)

type FollowService struct {
	followRepo  ports.FollowRepository
	userRepo    ports.UserRepository
	idGenerator *utils.IDGenerator
}

func NewFollowService(followRepo ports.FollowRepository, userRepo ports.UserRepository, idGenerator *utils.IDGenerator) *FollowService {
	return &FollowService{
		followRepo:  followRepo,
		userRepo:    userRepo,
		idGenerator: idGenerator,
	}
}

func (s *FollowService) FollowUser(followerID, followedID string) (*domain.Follow, error) {
	if followerID == "" || followedID == "" || followerID == followedID {
		return nil, errors.NewValidationError("invalid follow data")
	}

	follower, err := findUser(followerID, s.userRepo)
	if err != nil {
		return nil, err
	}

	followed, err := findUser(followedID, s.userRepo)
	if err != nil {
		return nil, err
	}

	existingFollow, err := s.followRepo.FindByFollowerAndFollowed(follower.ID, followed.ID)
	if err == nil && existingFollow != nil {
		return nil, errors.NewConflictError("user is already following this user")
	}

	id := s.idGenerator.Generate()
	follow := domain.NewFollow(id, followerID, followedID)

	if !follow.IsValid() {
		return nil, errors.NewValidationError("invalid follow data")
	}

	err = s.followRepo.Save(follow)
	if err != nil {
		return nil, errors.NewInternalError("failed to save follow")
	}

	return follow, nil
}

func (s *FollowService) UnfollowUser(followerID, followedID string) error {
	if followerID == "" || followedID == "" {
		return errors.NewValidationError("invalid follow data")
	}

	follow, err := s.followRepo.FindByFollowerAndFollowed(followerID, followedID)
	if err != nil {
		return errors.NewNotFoundError("follow relationship")
	}

	return s.followRepo.Delete(follow.ID)
}

func (s *FollowService) GetFollowing(userID string) ([]*domain.Follow, error) {
	if userID == "" {
		return nil, errors.NewValidationError("user ID cannot be empty")
	}

	user, err := findUser(userID, s.userRepo)
	if err != nil {
		return nil, err
	}

	return s.followRepo.FindFollowing(user.ID)
}

func (s *FollowService) GetFollowers(userID string) ([]*domain.Follow, error) {
	if userID == "" {
		return nil, errors.NewValidationError("user ID cannot be empty")
	}

	user, err := findUser(userID, s.userRepo)
	if err != nil {
		return nil, err
	}

	return s.followRepo.FindFollowers(user.ID)
}

func (s *FollowService) IsFollowing(followerID, followeeID string) (bool, error) {
	if followerID == "" || followeeID == "" {
		return false, errors.NewValidationError("invalid follow data")
	}

	_, err := s.followRepo.FindByFollowerAndFollowed(followerID, followeeID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// -----------------
//
//	Helper functions
//
// -----------------

func findUser(userID string, userRepo ports.UserRepository) (*domain.User, error) {
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}
	return user, nil
}
