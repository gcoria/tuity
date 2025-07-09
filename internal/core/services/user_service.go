package services

import (
	"tuity/internal/core/domain"
	"tuity/internal/core/ports"
	"tuity/pkg/errors"
	"tuity/pkg/utils"
)

type UserService struct {
	userRepo    ports.UserRepository
	idGenerator *utils.IDGenerator
}

func NewUserService(userRepo ports.UserRepository, idGenerator *utils.IDGenerator) *UserService {
	return &UserService{
		userRepo:    userRepo,
		idGenerator: idGenerator,
	}
}

func (s *UserService) CreateUser(username, displayName string) (*domain.User, error) {
	if username == "" || displayName == "" {
		return nil, errors.NewValidationError("username and display name cannot be empty")
	}
	existingUser, err := s.userRepo.FindByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.NewConflictError("username already exists")
	}

	id := s.idGenerator.Generate()
	user := domain.NewUser(id, username, displayName)
	if !user.IsValid() {
		return nil, errors.NewValidationError("invalid user data")
	}

	err = s.userRepo.Save(user)
	if err != nil {
		return nil, errors.NewInternalError("failed to save user")
	}

	return user, nil
}

func (s *UserService) GetUser(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.NewValidationError("user ID cannot be empty")
	}

	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	if username == "" {
		return nil, errors.NewValidationError("username cannot be empty")
	}

	return s.userRepo.FindByUsername(username)
}
