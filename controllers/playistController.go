package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tushar27x/music-lib-api/config"
	"github.com/tushar27x/music-lib-api/models"
	"gorm.io/gorm"
)

// @Summary     Add a new playlist
// @Description Create a new playlist with songs
// @Tags        playlists
// @Accept      json
// @Produce     json
// @Param       playlist body models.PlaylistCreateRequest true "Playlist data"
// @Success     200 {object} models.PlaylistResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /playlists/ [post]
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

// @Summary     Get all playlists
// @Description Retrieve all playlists for the authenticated user
// @Tags        playlists
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /playlists/ [get]
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

// @Summary     Get playlist by ID
// @Description Retrieve a specific playlist by its ID for the authenticated user
// @Tags        playlists
// @Produce     json
// @Param       id path int true "Playlist ID"
// @Success     200 {object} models.PlaylistResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /playlists/{id} [get]
func GetPlayListById(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	playlistId := c.Param("id")
	if playlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	var playlist models.Playlist
	if err := config.DB.Preload("Songs").
		Where("id = ? AND user_id = ?", playlistId, userId).
		First(&playlist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playlist": &playlist})
}

// @Summary     Update playlist by ID
// @Description Update a specific playlist by its ID for the authenticated user
// @Tags        playlists
// @Accept      json
// @Produce     json
// @Param       id path int true "Playlist ID"
// @Param       playlist body models.PlaylistCreateRequest true "Updated playlist data"
// @Success     200 {object} models.PlaylistResponse
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /playlists/{id} [put]
func UpdatePlaylist(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	playlistId := c.Param("id")
	if playlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	// Check if playlist exists and belongs to the user
	var existingPlaylist models.Playlist
	if err := config.DB.Where("id = ? AND user_id = ?", playlistId, userId).First(&existingPlaylist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the update data
	var updateData models.PlaylistCreateRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Update the playlist name
	if err := tx.Model(&existingPlaylist).Update("name", updateData.Name).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If song IDs are provided, update the playlist's songs
	if len(updateData.SongIds) > 0 {
		// Verify that all songs belong to the user
		var songs []models.Song
		if err := tx.Where("id IN ? AND user_id = ?", updateData.SongIds, userId).Find(&songs).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song IDs provided"})
			return
		}

		// Clear existing songs and add new ones
		if err := tx.Model(&existingPlaylist).Association("Songs").Clear(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing songs"})
			return
		}

		if err := tx.Model(&existingPlaylist).Association("Songs").Append(&songs); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add songs to playlist"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Return the updated playlist with songs
	if err := config.DB.Preload("Songs").First(&existingPlaylist, existingPlaylist.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"playlist": &existingPlaylist})
}

// @Summary     Delete playlist by ID
// @Description Soft delete a specific playlist by its ID (songs remain unaffected)
// @Tags        playlists
// @Produce     json
// @Param       id path int true "Playlist ID"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     404 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /playlists/{id} [delete]
func DeletePlaylist(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	playlistId := c.Param("id")
	if playlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	// Check if playlist exists and belongs to the user
	var playlist models.Playlist
	if err := config.DB.Where("id = ? AND user_id = ?", playlistId, userId).First(&playlist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()

	// Clear the association with songs (but don't delete the songs)
	if err := tx.Model(&playlist).Association("Songs").Clear(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear song associations"})
		return
	}

	// Soft delete the playlist
	if err := tx.Delete(&playlist).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted successfully (songs remain unaffected)"})
}

// @Summary     Search playlists
// @Description Search playlists by name with fuzzy matching
// @Tags        playlists
// @Produce     json
// @Param       q query string false "Search query (searches playlist name)"
// @Param       name query string false "Search by playlist name"
// @Param       limit query int false "Limit results (default: 20, max: 100)"
// @Param       offset query int false "Offset for pagination (default: 0)"
// @Success     200 {object} map[string]interface{}
// @Failure     400 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Security    BearerAuth
// @Router      /playlists/search [get]
func SearchPlaylists(c *gin.Context) {
	userId, ok := c.MustGet("userId").(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get search parameters
	query := c.Query("q")
	name := c.Query("name")
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
		// General search across playlist name
		dbQuery = dbQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+query+"%")
	} else if name != "" {
		// Specific name search
		dbQuery = dbQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%")
	}

	// Preload songs and apply pagination
	var playlists []models.Playlist
	var total int64

	// Count total results
	if err := dbQuery.Model(&models.Playlist{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get paginated results with songs
	if err := dbQuery.Preload("Songs").Limit(limit).Offset(offset).Find(&playlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"playlists": playlists,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_more": offset+limit < int(total),
		},
	})
}
