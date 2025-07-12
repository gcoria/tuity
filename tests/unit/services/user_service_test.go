package services

import (
	"testing"
	"tuity/internal/core/services"
	"tuity/pkg/errors"
	"tuity/pkg/utils"
	"tuity/tests/testutils/mocks"
)

func setupUserService() (*services.UserService, *mocks.MockUserRepository) {
	mockRepo := mocks.NewMockUserRepository()
	idGenerator := utils.NewIDGenerator()
	userService := services.NewUserService(mockRepo, idGenerator)
	return userService, mockRepo
}

func Test_CreateUser(t *testing.T) {
	userService, _ := setupUserService()

	tests := []struct {
		name        string
		username    string
		displayName string
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name:        "Valid user creation",
			username:    "diego",
			displayName: "Diego Maradona",
			expectError: false,
		},
		{
			name:        "Empty username",
			username:    "",
			displayName: "Diego Maradona",
			expectError: true,
			errorType:   errors.ValidationError,
		},
		{
			name:        "Empty display name",
			username:    "diego",
			displayName: "",
			expectError: true,
			errorType:   errors.ValidationError,
		},
		{
			name:        "Both empty",
			username:    "",
			displayName: "",
			expectError: true,
			errorType:   errors.ValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.CreateUser(tt.username, tt.displayName)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if domainErr, ok := err.(*errors.DomainError); ok {
					if domainErr.Type != tt.errorType {
						t.Errorf("Expected error type %s, got %s", tt.errorType, domainErr.Type)
					}
				} else {
					t.Error("Expected domain error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if user == nil {
					t.Error("Expected user to be created")
				}
				if user.Username != tt.username {
					t.Errorf("Expected username %s, got %s", tt.username, user.Username)
				}
				if user.DisplayName != tt.displayName {
					t.Errorf("Expected display name %s, got %s", tt.displayName, user.DisplayName)
				}
				if user.ID == "" {
					t.Error("Expected user ID to be set")
				}
			}
		})
	}
}

func Test_DuplicateUsername(t *testing.T) {
	userService, _ := setupUserService()

	_, err := userService.CreateUser("diego", "Diego Maradona")
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	_, err = userService.CreateUser("diego", "Diego Torres")
	if err == nil {
		t.Error("Expected error for duplicate username")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected domain error")
	}
	if domainErr.Type != errors.ConflictError {
		t.Errorf("Expected conflict error, got %s", domainErr.Type)
	}
}

func Test_RepositoryFailure(t *testing.T) {
	userService, mockRepo := setupUserService()
	mockRepo.ShouldFailSave = true

	_, err := userService.CreateUser("diego", "Diego Maradona")
	if err == nil {
		t.Error("Expected error for repository failure")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Error("Expected domain error")
	}
	if domainErr.Type != errors.InternalError {
		t.Errorf("Expected internal error, got %s", domainErr.Type)
	}
}

func Test_GetUser(t *testing.T) {
	userService, _ := setupUserService()

	createdUser, err := userService.CreateUser("diego", "Diego Maradona")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name        string
		userID      string
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name:        "Valid user ID",
			userID:      createdUser.ID,
			expectError: false,
		},
		{
			name:        "Empty user ID",
			userID:      "",
			expectError: true,
			errorType:   errors.ValidationError,
		},
		{
			name:        "Non-existent user ID",
			userID:      "non-existent",
			expectError: true,
			errorType:   errors.NotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.GetUser(tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if domainErr, ok := err.(*errors.DomainError); ok {
					if domainErr.Type != tt.errorType {
						t.Errorf("Expected error type %s, got %s", tt.errorType, domainErr.Type)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if user == nil {
					t.Error("Expected user to be found")
				}
				if user.ID != tt.userID {
					t.Errorf("Expected user ID %s, got %s", tt.userID, user.ID)
				}
			}
		})
	}
}

func Test_GetUserByUsername(t *testing.T) {
	userService, _ := setupUserService()

	createdUser, err := userService.CreateUser("diego", "Diego Maradona")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	tests := []struct {
		name        string
		username    string
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name:        "Valid username",
			username:    "diego",
			expectError: false,
		},
		{
			name:        "Empty username",
			username:    "",
			expectError: true,
			errorType:   errors.ValidationError,
		},
		{
			name:        "Non-existent username",
			username:    "pele",
			expectError: true,
			errorType:   errors.NotFoundError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.GetUserByUsername(tt.username)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if domainErr, ok := err.(*errors.DomainError); ok {
					if domainErr.Type != tt.errorType {
						t.Errorf("Expected error type %s, got %s", tt.errorType, domainErr.Type)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if user == nil {
					t.Error("Expected user to be found")
				}
				if user.Username != tt.username {
					t.Errorf("Expected username %s, got %s", tt.username, user.Username)
				}
				if user.ID != createdUser.ID {
					t.Errorf("Expected user ID %s, got %s", createdUser.ID, user.ID)
				}
			}
		})
	}
}
