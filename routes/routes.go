package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/controllers"
	"github.com/tushar27x/music-lib-api/middlewares"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		albums := api.Group("/albums")
		albums.Use(middlewares.AuthMiddleware())
		{
			albums.GET("/getAllAlbums", controllers.GetAlbums)
			albums.POST("/addAlbum", controllers.CreateAlbum)
		}

		songs := api.Group("/songs")
		songs.Use(middlewares.AuthMiddleware())
		{
			songs.GET("/getAllSongs", controllers.GetSongs)
			songs.POST("/addSong", controllers.AddSong)
		}

		playlists := api.Group("/playlists")
		playlists.Use(middlewares.AuthMiddleware())
		{
			playlists.GET("/getAllPlaylists", controllers.GetPlayList)
			playlists.POST("/addPlaylist", controllers.AddPlaylist)
		}
	}
}
