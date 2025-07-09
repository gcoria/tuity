package ports

import "tuity/internal/core/domain"

type TimelineRepository interface {
	Save(timeline *domain.Timeline) error
	FindByUserID(userID string) (*domain.Timeline, error)
	Delete(userID string) error
	Clear() error
}
