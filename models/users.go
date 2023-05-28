package models

import (
	"database/sql"
	"time"
)

// User is type for implement user model.
type User struct {
	ID         int          `gorm:"primaryKey" json:"id"`
	FirstName  string       `gorm:"size:50" json:"first_name"`
	LastName   string       `gorm:"size:50" json:"last_name"`
	Username   string       `gorm:"30;unique" json:"username"`
	Email      string       `gorm:"size:50;unique" json:"email"`
	Password   string       `gorm:"size:128" json:"-"`
	IsActive   bool         `gorm:"default:true" json:"is_active"`
	IsAdmin    bool         `gorm:"default:false" json:"is_admin"`
	LastLogin  sql.NullTime `gorm:"null" json:"last_login"`
	DateJoined time.Time    `gorm:"autoCreateTime" json:"date_joined"`
	UserBio    *UserBio     `gorm:"foreignKey:UserID" json:"user_bio"`
}

// UserBio is type for implament bio.
type UserBio struct {
	ID           int          `gorm:"primaryKey" json:"id"`
	UserID       int          `json:"user_id"`
	Avatar       string       `gorm:"size:255" json:"avatar"`
	IsAdmin      bool         `gorm:"default:false" json:"is_admin"`
	Gender       string       `gorm:"size:9" json:"gender"`
	StudentID    uint         `json:"student_id"`
	SchoolOrigin string       `gorm:"size:50" json:"school_origin"`
	Birth        sql.NullTime `json:"birth"`
	UpdatedAt    time.Time    `json:"updated_at"`
	CreatedAt    time.Time    `json:"created_at"`
}

// CreateNewUser is function to create a new user.
func CreateNewUser(user *User) error {
	user.Password = EncryptionPassword(user.Password)
	return DB().Create(user).Error
}

// GetUserByID is function to get user by id.
func GetUserByID(id int) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("id = ?", id).Preload("UserBio").First(&user).Error
	return user, err
}

// GetUserByEmail is function to get user by email.
func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("email = ?", email).Preload("UserBio").First(&user).Error
	return user, err
}

// GetUserByUsername is function to get user by username.
func GetUserByUsername(username string) (User, error) {
	var user User
	err := DB().Model(&User{}).Where("username = ?", username).Preload("UserBio").First(&user).Error
	return user, err
}
