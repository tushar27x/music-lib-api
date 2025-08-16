package models

import (
	"time"

	"gorm.io/gorm"
)

// Album represents a music album
// @Description Album model for organizing songs
type Album struct {
	// @Description Unique identifier for the album
	ID uint `json:"id" gorm:"primarykey" example:"1"`
	// @Description When the album was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the album was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the album was deleted (soft delete)
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	// @Description Album title
	Title string `json:"title" gorm:"unique" example:"Dark Side of the Moon"`
	// @Description Album artist
	Artist string `json:"artist" example:"Pink Floyd"`
	// @Description Release year
	Year int `json:"year" example:"1973"`
	// @Description User ID who owns the album
	UserId uint `json:"user_id" example:"1"`
	// @Description Songs in the album
	Songs []Song `json:"songs,omitempty" gorm:"foreignKey:AlbumId"`
}

// AlbumResponse represents the album data returned in API responses
// @Description Album response model
type AlbumResponse struct {
	// @Description Unique identifier for the album
	ID uint `json:"id" example:"1"`
	// @Description When the album was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the album was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description Album title
	Title string `json:"title" example:"Dark Side of the Moon"`
	// @Description Album artist
	Artist string `json:"artist" example:"Pink Floyd"`
	// @Description Release year
	Year int `json:"year" example:"1973"`
	// @Description User ID who owns the album
	UserId uint `json:"user_id" example:"1"`
	// @Description Songs in the album
	Songs []SongResponse `json:"songs,omitempty"`
}

// AlbumCreateRequest represents the album creation request payload
// @Description Album creation request model
type AlbumCreateRequest struct {
	// @Description Album title
	Title string `json:"title" binding:"required" example:"Dark Side of the Moon"`
	// @Description Album artist
	Artist string `json:"artist" binding:"required" example:"Pink Floyd"`
	// @Description Release year
	Year int `json:"year" binding:"required" example:"1973"`
}
