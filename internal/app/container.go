package app

import (
	"tuity/internal/adapters/driven/memory"
	"tuity/internal/adapters/driving/http/handlers"
	"tuity/internal/core/ports"
	"tuity/internal/core/services"
	"tuity/pkg/utils"
)

type Container struct {
	IDGenerator *utils.IDGenerator
	Config      *Config

	//Driven Adapters
	UserRepo     ports.UserRepository
	TweetRepo    ports.TweetRepository
	FollowRepo   ports.FollowRepository
	TimelineRepo ports.TimelineRepository
	CacheRepo    ports.CacheRepository

	// Core Business Logic
	UserService     *services.UserService
	TweetService    *services.TweetService
	FollowService   *services.FollowService
	TimelineService *services.TimelineService

	// Driving Adapters
	UserHandler     *handlers.UserHandler
	TweetHandler    *handlers.TweetHandler
	FollowHandler   *handlers.FollowHandler
	TimelineHandler *handlers.TimelineHandler
}

func NewContainer() *Container {
	config := LoadConfig()
	idGenerator := utils.NewIDGenerator()

	// Create repositories
	userRepo := memory.NewUserMemoryRepository()
	tweetRepo := memory.NewTweetMemoryRepository()
	followRepo := memory.NewFollowMemoryRepository()
	timelineRepo := memory.NewTimelineMemoryRepository()
	cacheRepo := memory.NewCacheMemoryRepository()

	// Create services
	userService := services.NewUserService(userRepo, idGenerator)
	tweetService := services.NewTweetService(tweetRepo, userRepo, idGenerator)
	followService := services.NewFollowService(followRepo, userRepo, idGenerator)
	timelineService := services.NewTimelineService(tweetRepo, followRepo, timelineRepo, userRepo, cacheRepo)

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	tweetHandler := handlers.NewTweetHandler(tweetService)
	followHandler := handlers.NewFollowHandler(followService)
	timelineHandler := handlers.NewTimelineHandler(timelineService, config.Timeline.DefaultLimit, config.Timeline.MaxLimit)

	return &Container{
		IDGenerator:     idGenerator,
		Config:          config,
		UserRepo:        userRepo,
		TweetRepo:       tweetRepo,
		FollowRepo:      followRepo,
		TimelineRepo:    timelineRepo,
		CacheRepo:       cacheRepo,
		UserService:     userService,
		TweetService:    tweetService,
		FollowService:   followService,
		TimelineService: timelineService,
		UserHandler:     userHandler,
		TweetHandler:    tweetHandler,
		FollowHandler:   followHandler,
		TimelineHandler: timelineHandler,
	}
}
