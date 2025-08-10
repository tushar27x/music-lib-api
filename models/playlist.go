package models

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint   `json:"userId"`
	Songs  []Song `json:"songs" gorm:"many2many:playlist_songs"`
}
