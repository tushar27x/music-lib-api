package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
	"gorm.io/gorm"
)

// @Summary     Add a new song
// @Description Add a new song to the user's library
// @Tags        songs
// @Accept      json
// @Produce     json
// @Param       song body models.SongCreateRequest true "Song data"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     403 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /songs/addSong [post]
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

// @Summary     Get all songs
// @Description Retrieve all songs for the authenticated user
// @Tags        songs
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /songs/getAllSongs [get]
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

// @Summary     Get song by ID
// @Description Retrieve a specific song by its ID for the authenticated user
// @Tags        songs
// @Produce     json
// @Param       id path int true "Song ID"
// @Success     200 {object} models.SongResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /songs/{id} [get]
func GetSongByID(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get song ID from URL parameter
	songId := c.Param("id")
	if songId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Song ID is required"})
		return
	}

	var song models.Song

	// Query the song, ensuring it belongs to the authenticated user
	if err := config.DB.Where("id = ? AND user_id = ?", songId, userId).First(&song).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"song": song})
}

// @Summary     Update song by ID
// @Description Update a specific song by its ID for the authenticated user
// @Tags        songs
// @Accept      json
// @Produce     json
// @Param       id path int true "Song ID"
// @Param       song body models.SongCreateRequest true "Updated song data"
// @Success     200 {object} models.SongResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	songId := c.Param("id")
	if songId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Song ID is required"})
		return
	}

	// Check if song exists and belongs to the user
	var existingSong models.Song
	if err := config.DB.Where("id = ? AND user_id = ?", songId, userId).First(&existingSong).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the update data
	var updateData models.SongCreateRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate album ownership if album_id is provided
	if updateData.AlbumId != nil {
		var album models.Album
		if err := config.DB.Where("id = ? AND user_id = ?", *updateData.AlbumId, userId).First(&album).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album_id or you don't own this album"})
			return
		}
	}

	// Update the song
	updates := map[string]interface{}{
		"title":    updateData.Title,
		"duration": updateData.Duration,
		"album_id": updateData.AlbumId,
	}

	if err := config.DB.Model(&existingSong).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated song
	c.JSON(http.StatusOK, gin.H{"song": existingSong})
}

// @Summary     Delete song by ID
// @Description Soft delete a specific song by its ID for the authenticated user
// @Tags        songs
// @Produce     json
// @Param       id path int true "Song ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	songId := c.Param("id")
	if songId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Song ID is required"})
		return
	}

	// Check if song exists and belongs to the user
	var song models.Song
	if err := config.DB.Where("id = ? AND user_id = ?", songId, userId).First(&song).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Remove song from all playlists (but don't delete the playlists)
	var playlists []models.Playlist
	if err := tx.Model(&song).Association("Playlists").Find(&playlists); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find associated playlists"})
		return
	}

	for _, playlist := range playlists {
		if err := tx.Model(&playlist).Association("Songs").Delete(&song); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove song from playlists"})
			return
		}
	}

	// Soft delete the song
	if err := tx.Delete(&song).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// @Summary     Search songs
// @Description Search songs by title, duration, or album with fuzzy matching
// @Tags        songs
// @Produce     json
// @Param       q query string false "Search query (searches title, duration)"
// @Param       title query string false "Search by title"
// @Param       album_id query int false "Search by album ID"
// @Param       min_duration query int false "Minimum duration in milliseconds"
// @Param       max_duration query int false "Maximum duration in milliseconds"
// @Param       limit query int false "Limit results (default: 20, max: 100)"
// @Param       offset query int false "Offset for pagination (default: 0)"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /songs/search [get]
func SearchSongs(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get search parameters
	query := c.Query("q")
	title := c.Query("title")
	albumIdStr := c.Query("album_id")
	minDurationStr := c.Query("min_duration")
	maxDurationStr := c.Query("max_duration")
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
		// General search across title and duration
		dbQuery = dbQuery.Where(
			"LOWER(title) LIKE LOWER(?) OR CAST(duration AS TEXT) LIKE ?",
			"%"+query+"%", "%"+query+"%",
		)
	} else {
		// Specific field searches
		if title != "" {
			dbQuery = dbQuery.Where("LOWER(title) LIKE LOWER(?)", "%"+title+"%")
		}
		if albumIdStr != "" {
			if albumId, err := strconv.Atoi(albumIdStr); err == nil {
				dbQuery = dbQuery.Where("album_id = ?", albumId)
			}
		}
		if minDurationStr != "" {
			if minDuration, err := strconv.Atoi(minDurationStr); err == nil {
				dbQuery = dbQuery.Where("duration >= ?", minDuration)
			}
		}
		if maxDurationStr != "" {
			if maxDuration, err := strconv.Atoi(maxDurationStr); err == nil {
				dbQuery = dbQuery.Where("duration <= ?", maxDuration)
			}
		}
	}

	// Preload album information and apply pagination
	var songs []models.Song
	var total int64

	// Count total results
	if err := dbQuery.Model(&models.Song{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get paginated results with album info
	if err := dbQuery.Preload("Album").Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songs": songs,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_more": offset+limit < int(total),
		},
	})
}
