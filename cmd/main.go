package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/routes"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	r := gin.Default()

	routes.RegisterRoutes(r)
	if err := r.Run(":" + config.GetEnv("PORT")); err != nil {
		log.Fatalf("Error occured while starting the server:%v", err)
	}
}
