package models

import (
	"database/sql"
	"time"
)

// User is type for implement user model.
type User struct {
	ID         int          `gorm:"primaryKey" json:"id,omitempty"`
	FirstName  string       `gorm:"size:50" json:"first_name,omitempty"`
	LastName   string       `gorm:"size:50" json:"last_name,omitempty"`
	Username   string       `gorm:"30;unique" json:"username,omitempty"`
	Email      string       `gorm:"size:50;unique" json:"email,omitempty"`
	Password   string       `gorm:"size:128" json:"-"`
	LastLogin  sql.NullTime `gorm:"null" json:"last_login,omitempty"`
	DateJoined time.Time    `gorm:"autoCreateTime" json:"date_joined,omitempty"`
	UserBio    *UserBio     `gorm:"foreignKey:UserID" json:"user_bio,omitempty"`
}

// UserBio is type for implament bio.
type UserBio struct {
	ID           int          `gorm:"primaryKey" json:"id,omitempty"`
	UserID       int          `json:"user_id,omitempty"`
	Avatar       string       `gorm:"size:255" json:"avatar,omitempty"`
	IsAdmin      bool         `gorm:"default:false" json:"is_admin,omitempty"`
	Gender       string       `gorm:"size:9" json:"gender,omitempty"`
	StudentID    uint         `json:"student_id,omitempty"`
	SchoolOrigin string       `gorm:"size:50" json:"school_origin,omitempty"`
	Birth        sql.NullTime `json:"birth,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
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
