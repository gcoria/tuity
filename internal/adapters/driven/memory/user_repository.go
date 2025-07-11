package memory

import (
	"sync"
	"tuity/internal/core/domain"
	"tuity/pkg/errors"
)

type UserMemoryRepository struct {
	users       map[string]*domain.User // Indexed by ID
	usersByName map[string]*domain.User // Indexed by username
	mutex       sync.RWMutex
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:       make(map[string]*domain.User),
		usersByName: make(map[string]*domain.User),
	}
}

func (r *UserMemoryRepository) Save(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID] = user
	r.usersByName[user.Username] = user
	return nil
}

func (r *UserMemoryRepository) FindByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.NewNotFoundError("user")
	}
	return user, nil
}

func (r *UserMemoryRepository) FindByUsername(username string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.usersByName[username]
	if !exists {
		return nil, errors.NewNotFoundError("user")
	}
	return user, nil
}

func (r *UserMemoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.NewNotFoundError("user")
	}

	delete(r.users, id)
	delete(r.usersByName, user.Username)
	return nil
}
