package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
// @Description User model for authentication and user management
type User struct {
	// @Description Unique identifier for the user
	ID uint `json:"id" gorm:"primarykey" example:"1"`
	// @Description When the user was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the user was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the user was deleted (soft delete)
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	// @Description User's full name
	Name string `json:"name" example:"John Doe"`
	// @Description User's email address (unique)
	Email string `json:"email" gorm:"unique" example:"john@example.com"`
	// @Description User's password (hashed)
	Password string `json:"password" example:"password123"`
	// @Description User's albums
	Albums []Album `json:"albums,omitempty" gorm:"foreignKey:UserId"`
	// @Description User's songs
	Songs []Song `json:"songs,omitempty" gorm:"foreignKey:UserId"`
	// @Description User's role in the system
	Role string `json:"role" example:"artist"`
}

// UserResponse represents the user data returned in API responses
// @Description User response model (excludes password)
type UserResponse struct {
	// @Description Unique identifier for the user
	ID uint `json:"id" example:"1"`
	// @Description When the user was created
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	// @Description When the user was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	// @Description User's full name
	Name string `json:"name" example:"John Doe"`
	// @Description User's email address
	Email string `json:"email" example:"john@example.com"`
	// @Description User's role in the system
	Role string `json:"role" example:"artist"`
	// @Description User's albums
	Albums []AlbumResponse `json:"albums,omitempty"`
	// @Description User's songs
	Songs []SongResponse `json:"songs,omitempty"`
}

// UserLoginRequest represents the login request payload
// @Description Login request model
type UserLoginRequest struct {
	// @Description User's email address
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	// @Description User's password
	Password string `json:"password" binding:"required" example:"password123"`
}

// UserRegisterRequest represents the registration request payload
// @Description Registration request model
type UserRegisterRequest struct {
	// @Description User's full name
	Name string `json:"name" binding:"required" example:"John Doe"`
	// @Description User's email address
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	// @Description User's password
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	// @Description User's role in the system
	Role string `json:"role" binding:"required" example:"artist"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
