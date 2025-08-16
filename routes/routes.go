package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tushar27x/music-lib-api/controllers"
	"github.com/tushar27x/music-lib-api/middlewares"
)

func RegisterRoutes(router *gin.Engine) {
	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			albums.GET("/", controllers.GetAlbums)
			albums.GET("/search", controllers.SearchAlbums)
			albums.GET("/:id", controllers.GetAlbumByID)
			albums.POST("/", controllers.CreateAlbum)
			albums.PUT("/:id", controllers.UpdateAlbum)
			albums.DELETE("/:id", controllers.DeleteAlbum)
		}

		songs := api.Group("/songs")
		songs.Use(middlewares.AuthMiddleware())
		{
			songs.GET("/", controllers.GetSongs)
			songs.GET("/search", controllers.SearchSongs)
			songs.GET("/:id", controllers.GetSongByID)
			songs.POST("/", controllers.AddSong)
			songs.PUT("/:id", controllers.UpdateSong)
			songs.DELETE("/:id", controllers.DeleteSong)
		}

		playlists := api.Group("/playlists")
		playlists.Use(middlewares.AuthMiddleware())
		{
			playlists.GET("/", controllers.GetPlayList)
			playlists.GET("/search", controllers.SearchPlaylists)
			playlists.GET("/:id", controllers.GetPlayListById)
			playlists.POST("/", controllers.AddPlaylist)
			playlists.PUT("/:id", controllers.UpdatePlaylist)
			playlists.DELETE("/:id", controllers.DeletePlaylist)
		}
	}
}
