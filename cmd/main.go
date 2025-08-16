// @title           Music Library API
// @version         1.0
// @description     A RESTful API for managing music library with albums, songs, and playlists
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      independent-carlene-tushar27x-a3461680.koyeb.app
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	_ "github.com/tushar27x/music-lib-api/docs"
	"github.com/tushar27x/music-lib-api/routes"
)

// @Summary     Health check endpoint
// @Description Get server health status
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /ping [get]
func main() {
	// Load environment variables
	config.LoadEnv()

	// Check for required environment variables
	requiredEnvVars := []string{"PORT", "DB_URI", "JWT_SECRET"}
	for _, envVar := range requiredEnvVars {
		if config.GetEnv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}

	// Connect to database
	config.ConnectDB()

	r := gin.Default()

	routes.RegisterRoutes(r)

	port := config.GetEnv("PORT")
	if port == "" {
		port = "8082" // fallback default
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error occurred while starting the server: %v", err)
	}
}
