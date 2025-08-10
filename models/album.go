package models

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Title  string `json:"title" gorm:"unique"`
	Artist string `json:"artist"`
	Year   int    `json:"year"`
	UserId uint   `json:"user_id"`
	Songs  []Song `json:"songs" gorm:"foreignKey:AlbumId"`
}
