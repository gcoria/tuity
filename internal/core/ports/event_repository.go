package ports

import "tuity/internal/core/domain"

type EventRepository interface {
	Save(event *domain.Event) error
	FindByID(id string) (*domain.Event, error)
	FindByUserID(userID string) ([]*domain.Event, error)
	FindByType(eventType string) ([]*domain.Event, error)
	FindRecent(limit int) ([]*domain.Event, error)
}
