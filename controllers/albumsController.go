package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
	"gorm.io/gorm"
)

func CreateAlbum(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	role, ok := c.MustGet("role").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if role != "artist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only artists can create albums"})
		return
	}
	album.UserId = userId

	if err := config.DB.Create(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, album)
}

func GetAlbums(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var albums []models.Album
	if err := config.DB.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}).Where("user_id = ?", userId).Find(&albums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, albums)
}
