package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique"`
	Password string  `json:"password"`
	Albums   []Album `json:"albums" gorm:"foreignKey:UserId"`
	Songs    []Song  `json:"songs" gorm:"foreignKey:UserId"`
	Role     string  `json:"role"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
