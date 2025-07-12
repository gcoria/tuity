package domain

import (
	"testing"
	"time"
)

func TestNewTweet(t *testing.T) {
	id := "tweet-123"
	userID := "user-123"
	content := "Hello world! This is my first tweet."

	tweet := NewTweet(id, userID, content)

	if tweet.ID != id {
		t.Errorf("Expected ID %s, got %s", id, tweet.ID)
	}
	if tweet.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, tweet.UserID)
	}
	if tweet.Content != content {
		t.Errorf("Expected Content %s, got %s", content, tweet.Content)
	}
	if tweet.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if tweet.IsDeleted {
		t.Error("Expected new tweet to not be deleted")
	}
}

func TestTweet_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		tweet         *Tweet
		expectedValid bool
	}{
		{
			name:          "Valid tweet",
			tweet:         NewTweet("tweet-123", "user-123", "Hello world!"),
			expectedValid: true,
		},
		{
			name:          "Empty ID",
			tweet:         NewTweet("", "user-123", "Hello world!"),
			expectedValid: false,
		},
		{
			name:          "Empty UserID",
			tweet:         NewTweet("tweet-123", "", "Hello world!"),
			expectedValid: false,
		},
		{
			name:          "Empty content",
			tweet:         NewTweet("tweet-123", "user-123", ""),
			expectedValid: false,
		},
		{
			name:          "Content at max length (280 chars)",
			tweet:         NewTweet("tweet-123", "user-123", generateString(280)),
			expectedValid: true,
		},
		{
			name:          "Content exceeds max length (281 chars)",
			tweet:         NewTweet("tweet-123", "user-123", generateString(281)),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tweet.IsValid() != tt.expectedValid {
				t.Errorf("Expected IsValid() to be %v, got %v", tt.expectedValid, tt.tweet.IsValid())
			}
		})
	}
}

func TestTweet_Delete(t *testing.T) {
	tweet := NewTweet("tweet-123", "user-123", "Hello world!")

	if tweet.IsDeleted {
		t.Error("Expected new tweet to not be deleted")
	}

	tweet.Delete()

	if !tweet.IsDeleted {
		t.Error("Expected tweet to be deleted after calling Delete()")
	}
}

func TestMaxTweetLength(t *testing.T) {
	expectedLength := 280
	if MaxTweetLength != expectedLength {
		t.Errorf("Expected MaxTweetLength to be %d, got %d", expectedLength, MaxTweetLength)
	}
}

func TestTweet_CreatedAtIsSet(t *testing.T) {
	beforeCreation := time.Now()
	tweet := NewTweet("tweet-123", "user-123", "Hello world!")
	afterCreation := time.Now()

	if tweet.CreatedAt.Before(beforeCreation) {
		t.Error("Tweet CreatedAt should be after creation started")
	}
	if tweet.CreatedAt.After(afterCreation) {
		t.Error("Tweet CreatedAt should be before creation finished")
	}
}

func generateString(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += "a"
	}
	return result
}
