package ports

import "time"

type CacheRepository interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
	Delete(key string)
	Flush()
	ItemCount() int
}
