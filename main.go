package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"MyHSRTrackerAPI/database"
	"MyHSRTrackerAPI/handlers"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using defaults")
	}

	// Initialize Database
	database.ConnectDB()

	router := gin.Default()

	// Phase 4: CORS and Security Middleware
	config := cors.DefaultConfig()
	origins := os.Getenv("ALLOWED_ORIGINS")
	config.AllowOrigins = strings.Split(origins, ",")
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// API Routes (Phase 3)
	api := router.Group("/api")
	{
		warp := api.Group("/warp")
		{
			warp.POST("/import", handlers.ImportWarp)
			warp.GET("/list", handlers.GetWarpList)
			warp.GET("/stats", handlers.GetWarpStats)
		}
	}

	port := os.Getenv("PORT")
	log.Printf("Server is running on http://localhost:%s\n", port)
	router.Run(":" + port)
}
