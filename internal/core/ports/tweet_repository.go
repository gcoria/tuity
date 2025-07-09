package ports

import "tuity/internal/core/domain"

type TweetRepository interface {
	Save(tweet *domain.Tweet) error
	FindByID(id string) (*domain.Tweet, error)
	FindByUserID(userID string) ([]*domain.Tweet, error)
	Delete(id string) error
}
