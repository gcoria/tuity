package app

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Cache     CacheConfig
	Timeline  TimelineConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type CacheConfig struct {
	TimelineTTL time.Duration
}

type TimelineConfig struct {
	DefaultLimit int
	MaxLimit     int
}

type RateLimitConfig struct {
	TweetCreateLimit     int
	FollowOperationLimit int
	TimelineRequestLimit int
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "release"),
		},
		Cache: CacheConfig{
			TimelineTTL: getEnvDuration("CACHE_TIMELINE_TTL", 5*time.Minute),
		},
		Timeline: TimelineConfig{
			DefaultLimit: getEnvInt("TIMELINE_DEFAULT_LIMIT", 20),
			MaxLimit:     getEnvInt("TIMELINE_MAX_LIMIT", 100),
		},
		RateLimit: RateLimitConfig{
			TweetCreateLimit:     getEnvInt("RATE_LIMIT_TWEET_CREATE", 10),
			FollowOperationLimit: getEnvInt("RATE_LIMIT_FOLLOW_OPS", 20),
			TimelineRequestLimit: getEnvInt("RATE_LIMIT_TIMELINE_REQUESTS", 100),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
