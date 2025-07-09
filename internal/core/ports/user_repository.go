package ports

import "tuity/internal/core/domain"

type UserRepository interface {
	Save(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Delete(id string) error
}
