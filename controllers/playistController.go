package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
)

func AddPlaylist(c *gin.Context) {
	var input models.Playlist
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userId").(uint)

	// Extract song IDs
	songIDs := make([]uint, 0, len(input.Songs))
	for _, song := range input.Songs {
		songIDs = append(songIDs, song.ID)
	}

	// Fetch only songs that belong to the user
	var songs []models.Song
	if len(songIDs) > 0 {
		if err := config.DB.
			Where("id IN ?", songIDs).
			Find(&songs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Create playlist
	playlist := models.Playlist{
		Name:   input.Name,
		UserId: userId,
	}

	if err := config.DB.Create(&playlist).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(songs) > 0 {
		if err := config.DB.Model(&playlist).Association("Songs").Append(&songs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Optionally return playlist with preloaded songs
	if err := config.DB.Preload("Songs").First(&playlist, playlist.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, playlist)
}

func GetPlayList(c *gin.Context) {
	var playlists []models.Playlist
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := config.DB.
		Preload("Songs").
		Where("user_id = ?", userId).
		Find(&playlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playlists": &playlists})
}
