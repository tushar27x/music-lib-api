package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Title    string `json:"title"`
	Duration uint   `json:"duration"`
	AlbumId  *uint  `json:"album_id"`
	UserId   uint   `json:"user_id"`
}
