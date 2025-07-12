package app

import (
	"tuity/internal/adapters/driving/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *Container) *gin.Engine {
	config := LoadConfig()
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(middleware.ErrorHandler())

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", container.UserHandler.CreateUser)
			users.GET("/:id", container.UserHandler.GetUser)
			users.GET("/username/:username", container.UserHandler.GetUserByUsername)

			users.GET("/:id/tweets", container.TweetHandler.GetUserTweets)

			timelineGroup := users.Group("")
			timelineGroup.Use(middleware.TimelineRequestRateLimit(config.RateLimit.TimelineRequestLimit))
			{
				timelineGroup.GET("/:id/timeline", container.TimelineHandler.GetTimeline)
				timelineGroup.POST("/:id/timeline/refresh", container.TimelineHandler.RefreshTimeline)
			}

			followGroup := users.Group("")
			followGroup.Use(middleware.FollowOperationRateLimit(config.RateLimit.FollowOperationLimit))
			{
				followGroup.POST("/:id/follow", container.FollowHandler.FollowUser)
				followGroup.DELETE("/:id/follow", container.FollowHandler.UnfollowUser)
			}

			users.GET("/:id/following", container.FollowHandler.GetFollowing)
			users.GET("/:id/followers", container.FollowHandler.GetFollowers)
			users.GET("/:id/following/:targetId", container.FollowHandler.IsFollowing)
		}

		tweets := v1.Group("/tweets")
		{
			tweets.POST("", middleware.TweetCreateRateLimit(config.RateLimit.TweetCreateLimit), container.TweetHandler.CreateTweet)
			tweets.GET("/:id", container.TweetHandler.GetTweet)
			tweets.DELETE("/:id", container.TweetHandler.DeleteTweet)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "tuity-api",
			"version": "1.0.0",
		})
	})

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
