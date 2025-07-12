package services

import (
	"sort"
	"time"
	"tuity/internal/core/domain"
	"tuity/internal/core/ports"
	"tuity/pkg/errors"
)

type TimelineService struct {
	tweetRepo    ports.TweetRepository
	followRepo   ports.FollowRepository
	timelineRepo ports.TimelineRepository
	userRepo     ports.UserRepository
	cache        ports.CacheRepository
}

func NewTimelineService(
	tweetRepo ports.TweetRepository,
	followRepo ports.FollowRepository,
	timelineRepo ports.TimelineRepository,
	userRepo ports.UserRepository,
	cache ports.CacheRepository,
) *TimelineService {
	return &TimelineService{
		tweetRepo:    tweetRepo,
		followRepo:   followRepo,
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
		cache:        cache,
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

	// Check cache first
	if cached, found := s.cache.Get(userID); found {
		if timeline, ok := cached.(*domain.Timeline); ok {
			limitedTweets := timeline.GetTweets(limit)
			result := &domain.Timeline{
				UserID:      timeline.UserID,
				Tweets:      limitedTweets,
				LastUpdated: timeline.LastUpdated,
			}
			return result, nil
		}
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	following, err := s.followRepo.FindFollowing(user.ID)
	if err != nil {
		return nil, errors.NewInternalError("failed to get following list")
	}

	userIDs := make([]string, 0, len(following)+1)
	userIDs = append(userIDs, userID)

	for _, follow := range following {
		userIDs = append(userIDs, follow.FollowedID)
	}

	allTweets := make([]*domain.Tweet, 0)
	for _, uid := range userIDs {
		userTweets, err := s.tweetRepo.FindByUserID(uid)
		if err == nil {
			for _, tweet := range userTweets {
				if !tweet.IsDeleted {
					allTweets = append(allTweets, tweet)
				}
			}
		}
	}

	sort.Slice(allTweets, func(i, j int) bool {
		return allTweets[i].CreatedAt.After(allTweets[j].CreatedAt)
	})

	timeline := domain.NewTimeline(userID)
	timeline.AddTweets(allTweets)

	s.cache.Set(userID, timeline, 0)

	limitedTweets := timeline.GetTweets(limit)
	timeline.Tweets = limitedTweets

	return timeline, nil
}

func (s *TimelineService) GetTimeline(userID string, limit int) (*domain.Timeline, error) {
	return s.GenerateTimeline(userID, limit)
}

func (s *TimelineService) RefreshTimeline(userID string) error {
	if userID == "" {
		return errors.NewValidationError("invalid user data")
	}

	s.cache.Delete(userID)

	return s.timelineRepo.Delete(userID)
}

func (s *TimelineService) InvalidateCache(userID string) {
	s.cache.Delete(userID)
}

func (s *TimelineService) SetCacheExpiration(userID string, duration time.Duration) {
	if cached, found := s.cache.Get(userID); found {
		s.cache.Set(userID, cached, duration)
	}
}

func (s *TimelineService) GetCacheStats() int {
	return s.cache.ItemCount()
}

func (s *TimelineService) FlushCache() {
	s.cache.Flush()
}
