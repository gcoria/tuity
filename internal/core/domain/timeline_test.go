package domain

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimeline(t *testing.T) {
	userID := "user-123"
	timeline := NewTimeline(userID)

	if timeline.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, timeline.UserID)
	}
	if timeline.LastUpdated.IsZero() {
		t.Error("Expected LastUpdated to be set")
	}
	if timeline.Tweets == nil {
		t.Error("Expected Tweets to be initialized")
	}
	if len(timeline.Tweets) != 0 {
		t.Errorf("Expected empty tweets list, got %d tweets", len(timeline.Tweets))
	}
}

func TestTimeline_AddTweets(t *testing.T) {
	timeline := NewTimeline("user-123")

	now := time.Now()
	tweet1 := &Tweet{
		ID:        "tweet-1",
		UserID:    "user-123",
		Content:   "First tweet",
		CreatedAt: now.Add(-2 * time.Hour),
		IsDeleted: false,
	}
	tweet2 := &Tweet{
		ID:        "tweet-2",
		UserID:    "user-456",
		Content:   "Second tweet",
		CreatedAt: now.Add(-1 * time.Hour),
		IsDeleted: false,
	}
	tweet3 := &Tweet{
		ID:        "tweet-3",
		UserID:    "user-789",
		Content:   "Third tweet",
		CreatedAt: now,
		IsDeleted: false,
	}

	tweets := []*Tweet{tweet1, tweet2, tweet3}
	timeline.AddTweets(tweets)

	if len(timeline.Tweets) != 3 {
		t.Errorf("Expected 3 tweets, got %d", len(timeline.Tweets))
	}

	if timeline.Tweets[0].ID != "tweet-3" {
		t.Errorf("Expected first tweet to be tweet-3, got %s", timeline.Tweets[0].ID)
	}
	if timeline.Tweets[1].ID != "tweet-2" {
		t.Errorf("Expected second tweet to be tweet-2, got %s", timeline.Tweets[1].ID)
	}
	if timeline.Tweets[2].ID != "tweet-1" {
		t.Errorf("Expected third tweet to be tweet-1, got %s", timeline.Tweets[2].ID)
	}
}

func TestTimeline_GetTweets(t *testing.T) {
	timeline := NewTimeline("user-123")

	tweets := make([]*Tweet, 5)
	for i := 0; i < 5; i++ {
		tweets[i] = &Tweet{
			ID:        fmt.Sprintf("tweet-%d", i),
			UserID:    "user-123",
			Content:   fmt.Sprintf("Tweet %d", i),
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			IsDeleted: false,
		}
	}
	timeline.AddTweets(tweets)

	tests := []struct {
		name          string
		limit         int
		expectedCount int
	}{
		{
			name:          "Get all tweets (limit 10)",
			limit:         10,
			expectedCount: 5,
		},
		{
			name:          "Get first 3 tweets",
			limit:         3,
			expectedCount: 3,
		},
		{
			name:          "Get first 1 tweet",
			limit:         1,
			expectedCount: 1,
		},
		{
			name:          "Zero limit returns all tweets",
			limit:         0,
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timeline.GetTweets(tt.limit)
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d tweets, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestTimeline_GetTweetsOrder(t *testing.T) {
	timeline := NewTimeline("user-123")

	now := time.Now()
	oldTweet := &Tweet{
		ID:        "old-tweet",
		UserID:    "user-123",
		Content:   "Old tweet",
		CreatedAt: now.Add(-3 * time.Hour),
		IsDeleted: false,
	}
	newTweet := &Tweet{
		ID:        "new-tweet",
		UserID:    "user-123",
		Content:   "New tweet",
		CreatedAt: now,
		IsDeleted: false,
	}
	middleTweet := &Tweet{
		ID:        "middle-tweet",
		UserID:    "user-123",
		Content:   "Middle tweet",
		CreatedAt: now.Add(-1 * time.Hour),
		IsDeleted: false,
	}

	tweets := []*Tweet{oldTweet, newTweet, middleTweet}
	timeline.AddTweets(tweets)

	result := timeline.GetTweets(3)

	if result[0].ID != "new-tweet" {
		t.Errorf("Expected newest tweet first, got %s", result[0].ID)
	}
	if result[1].ID != "middle-tweet" {
		t.Errorf("Expected middle tweet second, got %s", result[1].ID)
	}
	if result[2].ID != "old-tweet" {
		t.Errorf("Expected oldest tweet last, got %s", result[2].ID)
	}
}

func TestTimeline_GetTweetsLimitBehavior(t *testing.T) {
	timeline := NewTimeline("user-123")

	tweets := make([]*Tweet, 3)
	for i := 0; i < 3; i++ {
		tweets[i] = &Tweet{
			ID:        fmt.Sprintf("tweet-%d", i),
			UserID:    "user-123",
			Content:   fmt.Sprintf("Tweet %d", i),
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			IsDeleted: false,
		}
	}
	timeline.AddTweets(tweets)

	result := timeline.GetTweets(10)
	if len(result) != 3 {
		t.Errorf("Expected 3 tweets when limit > available, got %d", len(result))
	}

	result = timeline.GetTweets(3)
	if len(result) != 3 {
		t.Errorf("Expected 3 tweets when limit = available, got %d", len(result))
	}

	result = timeline.GetTweets(1)
	if len(result) != 1 {
		t.Errorf("Expected 1 tweet when limit < available, got %d", len(result))
	}
}

func TestTimeline_LastUpdatedIsSet(t *testing.T) {
	beforeCreation := time.Now()
	timeline := NewTimeline("user-123")
	afterCreation := time.Now()

	if timeline.LastUpdated.Before(beforeCreation) {
		t.Error("Timeline LastUpdated should be after creation started")
	}
	if timeline.LastUpdated.After(afterCreation) {
		t.Error("Timeline LastUpdated should be before creation finished")
	}
}

func TestTimeline_LastUpdatedOnAddTweets(t *testing.T) {
	timeline := NewTimeline("user-123")
	originalTime := timeline.LastUpdated

	time.Sleep(1 * time.Millisecond)

	tweets := []*Tweet{
		{
			ID:        "tweet-1",
			UserID:    "user-123",
			Content:   "Test tweet",
			CreatedAt: time.Now(),
			IsDeleted: false,
		},
	}

	timeline.AddTweets(tweets)

	if !timeline.LastUpdated.After(originalTime) {
		t.Error("LastUpdated should be updated when adding tweets")
	}
}
