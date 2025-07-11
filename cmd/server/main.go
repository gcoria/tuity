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

	gin.SetMode(gin.ReleaseMode)

	fmt.Println("Initializing dependencies...")
	container := app.NewContainer()

	fmt.Println("Setting up routes...")
	router := app.SetupRouter(container)

	printEndpoints()

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("API URL: http://localhost:%s/api/v1\n", port)
	fmt.Printf("Health Check: http://localhost:%s/health\n", port)
	fmt.Println("=====================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func printEndpoints() {
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
}
