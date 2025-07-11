package memory

import (
	"sync"
	"tuity/internal/core/domain"
	"tuity/pkg/errors"
)

type TweetMemoryRepository struct {
	tweets       map[string]*domain.Tweet
	tweetsByUser map[string][]*domain.Tweet
	mutex        sync.RWMutex
}

func NewTweetMemoryRepository() *TweetMemoryRepository {
	return &TweetMemoryRepository{
		tweets:       make(map[string]*domain.Tweet),
		tweetsByUser: make(map[string][]*domain.Tweet),
	}
}

func (r *TweetMemoryRepository) Save(tweet *domain.Tweet) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.tweets[tweet.ID]
	if exists {
		r.tweets[tweet.ID] = tweet

		userTweets := r.tweetsByUser[tweet.UserID]
		for i, t := range userTweets {
			if t.ID == tweet.ID {
				userTweets[i] = tweet
				break
			}
		}
	} else {
		r.tweets[tweet.ID] = tweet
		r.tweetsByUser[tweet.UserID] = append(r.tweetsByUser[tweet.UserID], tweet)
	}

	return nil
}

func (r *TweetMemoryRepository) FindByID(id string) (*domain.Tweet, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tweet, exists := r.tweets[id]
	if !exists {
		return nil, errors.NewNotFoundError("tweet")
	}
	return tweet, nil
}

func (r *TweetMemoryRepository) FindByUserID(userID string) ([]*domain.Tweet, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tweets, exists := r.tweetsByUser[userID]
	if !exists {
		return []*domain.Tweet{}, nil
	}

	result := make([]*domain.Tweet, len(tweets))
	copy(result, tweets)
	return result, nil
}

func (r *TweetMemoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tweet, exists := r.tweets[id]
	if !exists {
		return errors.NewNotFoundError("tweet")
	}

	delete(r.tweets, id)

	userTweets := r.tweetsByUser[tweet.UserID]
	for i, t := range userTweets {
		if t.ID == id {
			r.tweetsByUser[tweet.UserID] = append(userTweets[:i], userTweets[i+1:]...)
			break
		}
	}

	return nil
}
