package models

import (
	"time"

	"gorm.io/gorm"
)

// Song represents a music track
// @Description Song model for individual music tracks
type Song struct {
	// @Description Unique identifier for the song
	ID uint `json:"id" gorm:"primarykey" example:"1"`
	// @Description When the song was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the song was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the song was deleted (soft delete)
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	// @Description Song title
	Title string `json:"title" example:"Bohemian Rhapsody"`
	// @Description Song duration in milliseconds
	Duration uint `json:"duration" example:"157467"`
	// @Description Optional album ID the song belongs to
	AlbumId *uint `json:"album_id,omitempty" example:"1"`
	// @Description User ID who owns the song
	UserId uint `json:"user_id" example:"1"`
}

// SongResponse represents the song data returned in API responses
// @Description Song response model
type SongResponse struct {
	// @Description Unique identifier for the song
	ID uint `json:"id" example:"1"`
	// @Description When the song was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the song was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description Song title
	Title string `json:"title" example:"Bohemian Rhapsody"`
	// @Description Song duration in milliseconds
	Duration uint `json:"duration" example:"175000"`
	// @Description Optional album ID the song belongs to
	AlbumId *uint `json:"album_id,omitempty" example:"1"`
	// @Description User ID who owns the song
	UserId uint `json:"user_id" example:"1"`
}

// SongCreateRequest represents the song creation request payload
// @Description Song creation request model
type SongCreateRequest struct {
	// @Description Song title
	Title string `json:"title" binding:"required" example:"Bohemian Rhapsody"`
	// @Description Song duration in millseconds
	Duration uint `json:"duration" binding:"required" example:"175000"`
	// @Description Optional album ID the song belongs to
	AlbumId *uint `json:"album_id,omitempty" example:"1"`
}
