package services

import (
	"tuity/internal/core/domain"
	"tuity/internal/core/ports"
	"tuity/pkg/errors"
)

type TimelineService struct {
	tweetRepo    ports.TweetRepository
	followRepo   ports.FollowRepository
	timelineRepo ports.TimelineRepository
	userRepo     ports.UserRepository
}

func NewTimelineService(
	tweetRepo ports.TweetRepository,
	followRepo ports.FollowRepository,
	timelineRepo ports.TimelineRepository,
	userRepo ports.UserRepository,
) *TimelineService {
	return &TimelineService{
		tweetRepo:    tweetRepo,
		followRepo:   followRepo,
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
	}
}

const defaultLimit = 20

func (s *TimelineService) GenerateTimeline(userID string, limit int) (*domain.Timeline, error) {
	if userID == "" {
		return nil, errors.NewValidationError("invalid user data")
	}

	if limit <= 0 {
		limit = defaultLimit
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	following, err := s.followRepo.FindFollowing(user.ID)
	if err != nil {
		return nil, errors.NewInternalError("failed to get following list")
	}

	timeline := domain.NewTimeline(userID)
	allTweets := make([]*domain.Tweet, 0)

	userTweets, err := s.tweetRepo.FindByUserID(userID)
	if err == nil {
		for _, tweet := range userTweets {
			if !tweet.IsDeleted {
				allTweets = append(allTweets, tweet)
			}
		}
	}

	for _, follow := range following {
		followeeTweets, err := s.tweetRepo.FindByUserID(follow.FollowedID)
		if err == nil {
			for _, tweet := range followeeTweets {
				if !tweet.IsDeleted {
					allTweets = append(allTweets, tweet)
				}
			}
		}
	}

	timeline.AddTweets(allTweets)

	limitedTweets := timeline.GetTweets(limit)
	timeline.Tweets = limitedTweets

	return timeline, nil
}

func (s *TimelineService) GetTimeline(userID string, limit int) (*domain.Timeline, error) {
	// TODO: add in memory caching if have time
	return s.GenerateTimeline(userID, limit)
}

func (s *TimelineService) RefreshTimeline(userID string) error {
	if userID == "" {
		return errors.NewValidationError("invalid user data")
	}
	// TODO: clear cache
	return s.timelineRepo.Delete(userID)
}
