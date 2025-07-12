package testutils

import (
	"testing"
	"tuity/internal/core/domain"
	"tuity/pkg/errors"
)

func AssertError(t *testing.T, err error, expectedType errors.ErrorType) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got none")
		return
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected domain error")
		return
	}

	if domainErr.Type != expectedType {
		t.Errorf("Expected error type %s, got %s", expectedType, domainErr.Type)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func AssertValidTweet(t *testing.T, tweet *domain.Tweet, expectedUserID, expectedContent string) {
	t.Helper()
	if tweet == nil {
		t.Error("Expected tweet to be created")
		return
	}
	if tweet.UserID != expectedUserID {
		t.Errorf("Expected userID %s, got %s", expectedUserID, tweet.UserID)
	}
	if tweet.Content != expectedContent {
		t.Errorf("Expected content %s, got %s", expectedContent, tweet.Content)
	}
	if tweet.ID == "" {
		t.Error("Expected tweet ID to be set")
	}
}

func RepeatChar(char byte, count int) string {
	result := make([]byte, count)
	for i := 0; i < count; i++ {
		result[i] = char
	}
	return string(result)
}
