package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
	"gorm.io/gorm"
)

// @Summary     Create a new album
// @Description Create a new album (artists only)
// @Tags        albums
// @Accept      json
// @Produce     json
// @Param       album body models.AlbumCreateRequest true "Album data"
// @Success     201 {object} models.AlbumResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     403 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /albums/ [post]
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

// @Summary     Get all albums
// @Description Retrieve all albums for the authenticated user
// @Tags        albums
// @Produce     json
// @Success     200 {array} models.AlbumResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /albums/ [get]
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

// @Summary     Get album by ID
// @Description Retrieve a specific album by its ID for the authenticated user
// @Tags        albums
// @Produce     json
// @Param       id path int true "Album ID"
// @Success     200 {object} models.AlbumResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /albums/{id} [get]
func GetAlbumByID(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	// Get album ID from URL parameter
	albumId := c.Param("id")
	if albumId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID is required"})
		return
	}

	var album models.Album

	if err := config.DB.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}).Where("id = ? AND user_id = ?", albumId, userId).First(&album).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

// @Summary     Update album by ID
// @Description Update a specific album by its ID (artists only)
// @Tags        albums
// @Accept      json
// @Produce     json
// @Param       id path int true "Album ID"
// @Param       album body models.AlbumCreateRequest true "Updated album data"
// @Success     200 {object} models.AlbumResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     403 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /albums/{id} [put]
func UpdateAlbum(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	role, ok := c.MustGet("role").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user role"})
		return
	}

	if role != "artist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only artists can update albums"})
		return
	}

	albumId := c.Param("id")
	if albumId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID is required"})
		return
	}

	// Check if album exists and belongs to the user
	var existingAlbum models.Album
	if err := config.DB.Where("id = ? AND user_id = ?", albumId, userId).First(&existingAlbum).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the update data
	var updateData models.AlbumCreateRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the album
	updates := map[string]interface{}{
		"title":  updateData.Title,
		"artist": updateData.Artist,
		"year":   updateData.Year,
	}

	if err := config.DB.Model(&existingAlbum).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated album
	c.JSON(http.StatusOK, existingAlbum)
}

// @Summary     Delete album by ID
// @Description Delete a specific album by its ID and all its songs (artists only)
// @Tags        albums
// @Produce     json
// @Param       id path int true "Album ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     403 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /albums/{id} [delete]
func DeleteAlbum(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	role, ok := c.MustGet("role").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user role"})
		return
	}

	if role != "artist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only artists can delete albums"})
		return
	}

	albumId := c.Param("id")
	if albumId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID is required"})
		return
	}

	// Check if album exists and belongs to the user
	var album models.Album
	if err := config.DB.Where("id = ? AND user_id = ?", albumId, userId).First(&album).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Delete all songs associated with this album
	if err := tx.Where("album_id = ? AND user_id = ?", albumId, userId).Delete(&models.Song{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete songs: " + err.Error()})
		return
	}

	// Delete the album
	if err := tx.Delete(&album).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album and all its songs deleted successfully"})
}

// @Summary     Search albums
// @Description Search albums by title, artist, or year with fuzzy matching
// @Tags        albums
// @Produce     json
// @Param       q query string false "Search query (searches title, artist, year)"
// @Param       title query string false "Search by title"
// @Param       artist query string false "Search by artist"
// @Param       year query int false "Search by year"
// @Param       limit query int false "Limit results (default: 20, max: 100)"
// @Param       offset query int false "Offset for pagination (default: 0)"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /albums/search [get]
func SearchAlbums(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get search parameters
	query := c.Query("q")
	title := c.Query("title")
	artist := c.Query("artist")
	yearStr := c.Query("year")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	// Set default values
	limit := 20
	offset := 0

	// Parse limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				limit = 100
			} else {
				limit = parsedLimit
			}
		}
	}

	// Parse offset
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Build the query
	dbQuery := config.DB.Where("user_id = ?", userId)

	// Apply search filters
	if query != "" {
		// General search across title, artist, and year
		dbQuery = dbQuery.Where(
			"LOWER(title) LIKE LOWER(?) OR LOWER(artist) LIKE LOWER(?) OR CAST(year AS TEXT) LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%",
		)
	} else {
		// Specific field searches
		if title != "" {
			dbQuery = dbQuery.Where("LOWER(title) LIKE LOWER(?)", "%"+title+"%")
		}
		if artist != "" {
			dbQuery = dbQuery.Where("LOWER(artist) LIKE LOWER(?)", "%"+artist+"%")
		}
		if yearStr != "" {
			if year, err := strconv.Atoi(yearStr); err == nil {
				dbQuery = dbQuery.Where("year = ?", year)
			}
		}
	}

	// Preload songs and apply pagination
	var albums []models.Album
	var total int64

	// Count total results
	if err := dbQuery.Model(&models.Album{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get paginated results with songs
	if err := dbQuery.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}).Limit(limit).Offset(offset).Find(&albums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"albums": albums,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_more": offset+limit < int(total),
		},
	})
}
