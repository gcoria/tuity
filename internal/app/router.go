package app

import (
	"tuity/internal/adapters/driving/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *Container) *gin.Engine {
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

			users.POST("/:id/follow", container.FollowHandler.FollowUser)
			users.DELETE("/:id/follow", container.FollowHandler.UnfollowUser)

			users.GET("/:id/following", container.FollowHandler.GetFollowing)
			users.GET("/:id/followers", container.FollowHandler.GetFollowers)
			users.GET("/:id/following/:targetId", container.FollowHandler.IsFollowing)
		}

		tweets := v1.Group("/tweets")
		{
			tweets.POST("", container.TweetHandler.CreateTweet)
			tweets.GET("/:id", container.TweetHandler.GetTweet)
			tweets.DELETE("/:id", container.TweetHandler.DeleteTweet)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "tuity-api",
		})
	})

	return router
}

// corsMiddleware adds CORS headers for frontend development
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
