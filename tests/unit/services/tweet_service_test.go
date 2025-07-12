package services

import (
	"testing"
	"tuity/internal/core/domain"
	"tuity/internal/core/services"
	"tuity/pkg/errors"
	"tuity/pkg/utils"
	"tuity/tests/testutils"
	"tuity/tests/testutils/mocks"
)

const (
	testUserID   = "user-123"
	testUsername = "lio"
	testContent  = "Hello world!"
	maxLength    = 280
)

func setupService() (*services.TweetService, *mocks.MockTweetRepository, *mocks.MockUserRepository) {
	tweetRepo := mocks.NewMockTweetRepository()
	userRepo := mocks.NewMockUserRepository()
	idGen := utils.NewIDGenerator()
	service := services.NewTweetService(tweetRepo, userRepo, idGen)

	userRepo.Save(&domain.User{
		ID:          testUserID,
		Username:    testUsername,
		DisplayName: "Lio Messi",
	})

	return service, tweetRepo, userRepo
}

func TestTweetService_CreateTweet(t *testing.T) {
	service, _, _ := setupService()

	tests := []struct {
		name      string
		userID    string
		content   string
		wantError bool
		errorType errors.ErrorType
	}{
		{"valid tweet", testUserID, testContent, false, ""},
		{"empty userID", "", testContent, true, errors.ValidationError},
		{"empty content", testUserID, "", true, errors.ValidationError},
		{"max length content", testUserID, testutils.RepeatChar('a', maxLength), false, ""},
		{"too long content", testUserID, testutils.RepeatChar('a', maxLength+1), true, errors.ValidationError},
		{"nonexistent user", "bad-user", testContent, true, errors.NotFoundError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tweet, err := service.CreateTweet(tt.userID, tt.content)

			if tt.wantError {
				testutils.AssertError(t, err, tt.errorType)
			} else {
				testutils.AssertNoError(t, err)
				testutils.AssertValidTweet(t, tweet, tt.userID, tt.content)
			}
		})
	}
}

func TestTweetService_CreateTweet_RepositoryFailure(t *testing.T) {
	service, tweetRepo, _ := setupService()
	tweetRepo.ShouldFailSave = true

	_, err := service.CreateTweet(testUserID, testContent)
	testutils.AssertError(t, err, errors.InternalError)
}

func TestTweetService_GetTweet(t *testing.T) {
	service, _, _ := setupService()

	tweet, err := service.CreateTweet(testUserID, testContent)
	testutils.AssertNoError(t, err)

	tests := []struct {
		name      string
		tweetID   string
		wantError bool
		errorType errors.ErrorType
	}{
		{"valid ID", tweet.ID, false, ""},
		{"empty ID", "", true, errors.ValidationError},
		{"nonexistent ID", "bad-id", true, errors.NotFoundError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetTweet(tt.tweetID)

			if tt.wantError {
				testutils.AssertError(t, err, tt.errorType)
			} else {
				testutils.AssertNoError(t, err)
				if result.ID != tt.tweetID {
					t.Errorf("Expected tweet ID %s, got %s", tt.tweetID, result.ID)
				}
			}
		})
	}
}

func TestTweetService_GetTweet_DeletedTweet(t *testing.T) {
	service, tweetRepo, _ := setupService()

	deletedTweet := &domain.Tweet{
		ID:        "deleted-tweet",
		UserID:    testUserID,
		Content:   "Deleted content",
		IsDeleted: true,
	}
	tweetRepo.Save(deletedTweet)

	_, err := service.GetTweet("deleted-tweet")
	testutils.AssertError(t, err, errors.NotFoundError)
}

func TestTweetService_GetUserTweets(t *testing.T) {
	service, tweetRepo, _ := setupService()

	tweet1, _ := service.CreateTweet(testUserID, "First tweet")
	tweet2, _ := service.CreateTweet(testUserID, "Second tweet")

	deletedTweet := &domain.Tweet{
		ID:        "deleted-tweet",
		UserID:    testUserID,
		Content:   "Deleted content",
		IsDeleted: true,
	}
	tweetRepo.Save(deletedTweet)

	tweets, err := service.GetUserTweets(testUserID)
	testutils.AssertNoError(t, err)

	if len(tweets) != 2 {
		t.Errorf("Expected 2 tweets, got %d", len(tweets))
	}

	tweetIDs := make(map[string]bool)
	for _, tweet := range tweets {
		if tweet.IsDeleted {
			t.Error("Deleted tweet should not be included")
		}
		tweetIDs[tweet.ID] = true
	}

	if !tweetIDs[tweet1.ID] || !tweetIDs[tweet2.ID] {
		t.Error("Expected tweets not found in results")
	}
}

func TestTweetService_DeleteTweet(t *testing.T) {
	service, tweetRepo, _ := setupService()

	tweet, err := service.CreateTweet(testUserID, testContent)
	testutils.AssertNoError(t, err)

	tests := []struct {
		name      string
		tweetID   string
		userID    string
		wantError bool
		errorType errors.ErrorType
	}{
		{"valid deletion", tweet.ID, testUserID, false, ""},
		{"empty tweet ID", "", testUserID, true, errors.ValidationError},
		{"empty user ID", tweet.ID, "", true, errors.ValidationError},
		{"nonexistent tweet", "bad-id", testUserID, true, errors.NotFoundError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteTweet(tt.tweetID, tt.userID)

			if tt.wantError {
				testutils.AssertError(t, err, tt.errorType)
			} else {
				testutils.AssertNoError(t, err)
				deletedTweet, _ := tweetRepo.FindByID(tt.tweetID)
				if deletedTweet != nil && !deletedTweet.IsDeleted {
					t.Error("Tweet should be marked as deleted")
				}
			}
		})
	}
}
