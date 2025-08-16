// @title           Music Library API
// @version         1.0
// @description     A RESTful API for managing music library with albums, songs, and playlists
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      {{.PROD_HOST}}
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
	config.LoadEnv()
	config.ConnectDB()

	r := gin.Default()

	routes.RegisterRoutes(r)
	if err := r.Run(":" + config.GetEnv("PORT")); err != nil {
		log.Fatalf("Error occured while starting the server:%v", err)
	}
}
