package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	id := "user-123"
	username := "lio_messi"
	displayName := "Lionel Messi"

	user := NewUser(id, username, displayName)

	if user.ID != id {
		t.Errorf("Expected ID %s, got %s", id, user.ID)
	}
	if user.Username != username {
		t.Errorf("Expected Username %s, got %s", username, user.Username)
	}
	if user.DisplayName != displayName {
		t.Errorf("Expected DisplayName %s, got %s", displayName, user.DisplayName)
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		expectedValid bool
	}{
		{
			name:          "Valid user",
			user:          NewUser("user-123", "lio_messi", "Lionel Messi"),
			expectedValid: true,
		},
		{
			name:          "Empty ID",
			user:          NewUser("", "lio_messi", "Lionel Messi"),
			expectedValid: false,
		},
		{
			name:          "Empty username",
			user:          NewUser("user-123", "", "Lionel Messi"),
			expectedValid: false,
		},
		{
			name:          "Empty display name",
			user:          NewUser("user-123", "lio_messi", ""),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.IsValid() != tt.expectedValid {
				t.Errorf("Expected IsValid() to be %v, got %v", tt.expectedValid, tt.user.IsValid())
			}
		})
	}
}
