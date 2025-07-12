package mocks

import (
	"tuity/internal/core/domain"
	"tuity/pkg/errors"
)

// MockUserRepository implements the UserRepository interface for testing
type MockUserRepository struct {
	users          map[string]*domain.User
	usersByName    map[string]*domain.User
	ShouldFailSave bool
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:       make(map[string]*domain.User),
		usersByName: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Save(user *domain.User) error {
	if m.ShouldFailSave {
		return errors.NewInternalError("failed to save user")
	}
	m.users[user.ID] = user
	m.usersByName[user.Username] = user
	return nil
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.NewNotFoundError("user")
}

func (m *MockUserRepository) FindByUsername(username string) (*domain.User, error) {
	if user, exists := m.usersByName[username]; exists {
		return user, nil
	}
	return nil, errors.NewNotFoundError("user")
}

func (m *MockUserRepository) Delete(id string) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.usersByName, user.Username)
		return nil
	}
	return errors.NewNotFoundError("user")
}

// MockTweetRepository implements the TweetRepository interface for testing
type MockTweetRepository struct {
	tweets         map[string]*domain.Tweet
	tweetsByUser   map[string][]*domain.Tweet
	ShouldFailSave bool
}

// NewMockTweetRepository creates a new mock tweet repository
func NewMockTweetRepository() *MockTweetRepository {
	return &MockTweetRepository{
		tweets:       make(map[string]*domain.Tweet),
		tweetsByUser: make(map[string][]*domain.Tweet),
	}
}

func (m *MockTweetRepository) Save(tweet *domain.Tweet) error {
	if m.ShouldFailSave {
		return errors.NewInternalError("failed to save tweet")
	}
	m.tweets[tweet.ID] = tweet
	m.tweetsByUser[tweet.UserID] = append(m.tweetsByUser[tweet.UserID], tweet)
	return nil
}

func (m *MockTweetRepository) FindByID(id string) (*domain.Tweet, error) {
	if tweet, exists := m.tweets[id]; exists {
		return tweet, nil
	}
	return nil, errors.NewNotFoundError("tweet")
}

func (m *MockTweetRepository) FindByUserID(userID string) ([]*domain.Tweet, error) {
	if tweets, exists := m.tweetsByUser[userID]; exists {
		return tweets, nil
	}
	return []*domain.Tweet{}, nil
}

func (m *MockTweetRepository) Delete(id string) error {
	if tweet, exists := m.tweets[id]; exists {
		delete(m.tweets, id)
		userTweets := m.tweetsByUser[tweet.UserID]
		for i, t := range userTweets {
			if t.ID == id {
				m.tweetsByUser[tweet.UserID] = append(userTweets[:i], userTweets[i+1:]...)
				break
			}
		}
		return nil
	}
	return errors.NewNotFoundError("tweet")
}
