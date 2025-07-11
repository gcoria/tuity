package memory

import (
	"sync"
	"tuity/internal/core/domain"
	"tuity/pkg/errors"
)

type TimelineMemoryRepository struct {
	timelines map[string]*domain.Timeline
	mutex     sync.RWMutex
}

func NewTimelineMemoryRepository() *TimelineMemoryRepository {
	return &TimelineMemoryRepository{
		timelines: make(map[string]*domain.Timeline),
	}
}

func (r *TimelineMemoryRepository) Save(timeline *domain.Timeline) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.timelines[timeline.UserID] = timeline
	return nil
}

func (r *TimelineMemoryRepository) FindByUserID(userID string) (*domain.Timeline, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	timeline, exists := r.timelines[userID]
	if !exists {
		return nil, errors.NewNotFoundError("timeline")
	}
	return timeline, nil
}

func (r *TimelineMemoryRepository) Delete(userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.timelines[userID]; !exists {
		return errors.NewNotFoundError("timeline")
	}

	delete(r.timelines, userID)
	return nil
}

func (r *TimelineMemoryRepository) Clear() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.timelines = make(map[string]*domain.Timeline)
	return nil
}
