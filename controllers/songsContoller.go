package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
)

func AddSong(c *gin.Context) {
	var song models.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Optional album ownership check if AlbumId provided
	if song.AlbumId != nil {
		var album models.Album
		if err := config.DB.First(&album, *song.AlbumId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album_id"})
			return
		}
		if album.UserId != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this album"})
			return
		}
	}

	song.UserId = userId

	if err := config.DB.Create(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"song": song})
}

func GetSongs(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var songs []models.Song
	if err := config.DB.Where("user_id = ?", userId).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}
