package domain

import (
	"testing"
	"time"
)

func TestNewFollow(t *testing.T) {
	id := "follow-123"
	followerID := "user-123"
	followedID := "user-456"

	follow := NewFollow(id, followerID, followedID)

	if follow.ID != id {
		t.Errorf("Expected ID %s, got %s", id, follow.ID)
	}
	if follow.FollowerID != followerID {
		t.Errorf("Expected FollowerID %s, got %s", followerID, follow.FollowerID)
	}
	if follow.FollowedID != followedID {
		t.Errorf("Expected FollowedID %s, got %s", followedID, follow.FollowedID)
	}
	if follow.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestFollow_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		follow        *Follow
		expectedValid bool
	}{
		{
			name:          "Valid follow",
			follow:        NewFollow("follow-123", "user-123", "user-456"),
			expectedValid: true,
		},
		{
			name:          "Empty ID",
			follow:        NewFollow("", "user-123", "user-456"),
			expectedValid: false,
		},
		{
			name:          "Empty FollowerID",
			follow:        NewFollow("follow-123", "", "user-456"),
			expectedValid: false,
		},
		{
			name:          "Empty FollowedID",
			follow:        NewFollow("follow-123", "user-123", ""),
			expectedValid: false,
		},
		{
			name:          "Self follow (same follower and followed)",
			follow:        NewFollow("follow-123", "user-123", "user-123"),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.follow.IsValid() != tt.expectedValid {
				t.Errorf("Expected IsValid() to be %v, got %v", tt.expectedValid, tt.follow.IsValid())
			}
		})
	}
}

func TestFollow_CreatedAtIsSet(t *testing.T) {
	beforeCreation := time.Now()
	follow := NewFollow("follow-123", "user-123", "user-456")
	afterCreation := time.Now()

	if follow.CreatedAt.Before(beforeCreation) {
		t.Error("Follow CreatedAt should be after creation started")
	}
	if follow.CreatedAt.After(afterCreation) {
		t.Error("Follow CreatedAt should be before creation finished")
	}
}

func TestFollow_AsymmetricRelationship(t *testing.T) {
	followAB := NewFollow("follow-1", "user-A", "user-B")
	followBA := NewFollow("follow-2", "user-B", "user-A")

	if !followAB.IsValid() {
		t.Error("Follow A->B should be valid")
	}
	if !followBA.IsValid() {
		t.Error("Follow B->A should be valid")
	}

	if followAB.ID == followBA.ID {
		t.Error("Different follows should have different IDs")
	}
	if followAB.FollowerID == followBA.FollowerID {
		t.Error("Different follows should have different follower relationships")
	}
}
