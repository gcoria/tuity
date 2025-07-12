package memory

import (
	"sync"
	"time"
	"tuity/internal/core/ports"

	"github.com/patrickmn/go-cache"
)

type CacheMemoryRepository struct {
	cache *cache.Cache
	mutex sync.RWMutex
}

func NewCacheMemoryRepository() ports.CacheRepository {
	c := cache.New(10*time.Minute, 15*time.Minute)

	return &CacheMemoryRepository{
		cache: c,
	}
}

func (r *CacheMemoryRepository) Get(key string) (interface{}, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.cache.Get(key)
}

func (r *CacheMemoryRepository) Set(key string, value interface{}, expiration time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cache.Set(key, value, expiration)
}

func (r *CacheMemoryRepository) Delete(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cache.Delete(key)
}

func (r *CacheMemoryRepository) Flush() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cache.Flush()
}

func (r *CacheMemoryRepository) ItemCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.cache.ItemCount()
}
