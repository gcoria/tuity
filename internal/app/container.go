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

	//Driven Adapters
	UserRepo     ports.UserRepository
	TweetRepo    ports.TweetRepository
	FollowRepo   ports.FollowRepository
	TimelineRepo ports.TimelineRepository

	// Core Business Logic
	UserService     *services.UserService
	TweetService    *services.TweetService
	FollowService   *services.FollowService
	TimelineService *services.TimelineService

	// Driving Adapters
	UserHandler *handlers.UserHandler
}

func NewContainer() *Container {

	idGenerator := utils.NewIDGenerator()

	userRepo := memory.NewUserMemoryRepository()
	tweetRepo := memory.NewTweetMemoryRepository()
	followRepo := memory.NewFollowMemoryRepository()
	timelineRepo := memory.NewTimelineMemoryRepository()

	userService := services.NewUserService(userRepo, idGenerator)
	tweetService := services.NewTweetService(tweetRepo, userRepo, idGenerator)
	followService := services.NewFollowService(followRepo, userRepo, idGenerator)
	timelineService := services.NewTimelineService(tweetRepo, followRepo, timelineRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)

	return &Container{
		IDGenerator:     idGenerator,
		UserRepo:        userRepo,
		TweetRepo:       tweetRepo,
		FollowRepo:      followRepo,
		TimelineRepo:    timelineRepo,
		UserService:     userService,
		TweetService:    tweetService,
		FollowService:   followService,
		TimelineService: timelineService,
		UserHandler:     userHandler,
	}
}
