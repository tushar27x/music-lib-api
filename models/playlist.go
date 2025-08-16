package models

import (
	"time"

	"gorm.io/gorm"
)

// Playlist represents a collection of songs
// @Description Playlist model for organizing songs into collections
type Playlist struct {
	// @Description Unique identifier for the playlist
	ID uint `json:"id" gorm:"primarykey" example:"1"`
	// @Description When the playlist was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the playlist was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the playlist was deleted (soft delete)
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	// @Description Playlist name
	Name string `json:"name" example:"My Favorite Songs"`
	// @Description User ID who owns the playlist
	UserId uint `json:"userId" example:"1"`
	// @Description Songs in the playlist
	Songs []Song `json:"songs,omitempty" gorm:"many2many:playlist_songs"`
}

// PlaylistResponse represents the playlist data returned in API responses
// @Description Playlist response model
type PlaylistResponse struct {
	// @Description Unique identifier for the playlist
	ID uint `json:"id" example:"1"`
	// @Description When the playlist was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the playlist was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description Playlist name
	Name string `json:"name" example:"My Favorite Songs"`
	// @Description User ID who owns the playlist
	UserId uint `json:"userId" example:"1"`
	// @Description Songs in the playlist
	Songs []SongResponse `json:"songs,omitempty"`
}

// PlaylistCreateRequest represents the playlist creation request payload
// @Description Playlist creation request model
type PlaylistCreateRequest struct {
	// @Description Playlist name
	Name string `json:"name" binding:"required" example:"My Favorite Songs"`
	// @Description Array of song IDs to include in the playlist
	SongIds []uint `json:"song_ids,omitempty" example:"[1,2,3]"`
}
