package services

import (
	"tuity/internal/core/domain"
	"tuity/internal/core/ports"
	"tuity/pkg/errors"
	"tuity/pkg/utils"
)

type TweetService struct {
	tweetRepo   ports.TweetRepository
	userRepo    ports.UserRepository
	idGenerator *utils.IDGenerator
}

func NewTweetService(tweetRepo ports.TweetRepository, userRepo ports.UserRepository, idGenerator *utils.IDGenerator) *TweetService {
	return &TweetService{
		tweetRepo:   tweetRepo,
		userRepo:    userRepo,
		idGenerator: idGenerator,
	}
}

func (s *TweetService) CreateTweet(userID, content string) (*domain.Tweet, error) {
	if userID == "" || content == "" {
		return nil, errors.NewValidationError("invalid tweet data")
	}
	if len(content) > domain.MaxTweetLength {
		return nil, errors.NewValidationError("tweet content exceeds maximum length of 280 characters")
	}

	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}

	id := s.idGenerator.Generate()
	tweet := domain.NewTweet(id, userID, content)

	if !tweet.IsValid() {
		return nil, errors.NewValidationError("invalid tweet data")
	}

	err = s.tweetRepo.Save(tweet)
	if err != nil {
		return nil, errors.NewInternalError("failed to save tweet")
	}

	return tweet, nil
}

func (s *TweetService) GetTweet(id string) (*domain.Tweet, error) {
	if id == "" {
		return nil, errors.NewValidationError("tweet ID cannot be empty")
	}

	tweet, err := s.tweetRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if tweet.IsDeleted {
		return nil, errors.NewNotFoundError("tweet")
	}

	return tweet, nil
}

func (s *TweetService) GetUserTweets(userID string) ([]*domain.Tweet, error) {
	if userID == "" {
		return nil, errors.NewValidationError("invalid user data")
	}

	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}

	tweets, err := s.tweetRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	activeTweets := make([]*domain.Tweet, 0)
	for _, tweet := range tweets {
		if !tweet.IsDeleted {
			activeTweets = append(activeTweets, tweet)
		}
	}

	return activeTweets, nil
}

func (s *TweetService) DeleteTweet(id, userID string) error {
	if id == "" || userID == "" {
		return errors.NewValidationError("invalid tweet data")
	}

	tweet, err := s.tweetRepo.FindByID(id)
	if err != nil {
		return err
	}

	if tweet.UserID != userID {
		return errors.NewValidationError("user can only delete their own tweets")
	}

	if tweet.IsDeleted {
		return errors.NewNotFoundError("tweet")
	}

	tweet.Delete()
	return s.tweetRepo.Save(tweet)
}
