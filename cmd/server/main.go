package main

import (
	"fmt"
	"log"
	"tuity/internal/app"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Starting Tuity")
	fmt.Println("=====================================")

	// Load configuration
	config := app.LoadConfig()

	gin.SetMode(config.Server.Mode)

	fmt.Println("Initializing dependencies...")
	container := app.NewContainer()

	fmt.Println("Setting up routes...")
	router := app.SetupRouter(container)

	printEndpoints(config.Server.Port)

	fmt.Printf("Server starting on port %s...\n", config.Server.Port)
	fmt.Printf("API URL: http://localhost:%s/api/v1\n", config.Server.Port)
	fmt.Printf("Health Check: http://localhost:%s/health\n", config.Server.Port)
	fmt.Println("=====================================")

	if err := router.Run(":" + config.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func printEndpoints(port string) {
	fmt.Println("\nAvailable Endpoints:")
	fmt.Println("Health:")
	fmt.Println("  GET    /health")

	fmt.Println("\nUsers:")
	fmt.Println("  POST   /api/v1/users")
	fmt.Println("  GET    /api/v1/users/:id")
	fmt.Println("  GET    /api/v1/users/:username")

	fmt.Println("\nTweets:")
	fmt.Println("  POST   /api/v1/tweets")
	fmt.Println("  GET    /api/v1/tweets/:id")
	fmt.Println("  DELETE /api/v1/tweets/:id")
	fmt.Println("  GET    /api/v1/users/:id/tweets")

	fmt.Println("\nFollow:")
	fmt.Println("  POST   /api/v1/users/:id/follow")
	fmt.Println("  DELETE /api/v1/users/:id/follow")
	fmt.Println("  GET    /api/v1/users/:id/following")
	fmt.Println("  GET    /api/v1/users/:id/followers")
	fmt.Println("  GET    /api/v1/users/:id/following/:targetId")

	fmt.Println("\nTimeline:")
	fmt.Println("  GET    /api/v1/users/:id/timeline")
	fmt.Println("  POST   /api/v1/users/:id/timeline/refresh")

	fmt.Println("\n💡 Authentication: Include 'X-User-ID' header for protected endpoints")
	fmt.Printf("🔧 Environment: %s\n", gin.Mode())
}
