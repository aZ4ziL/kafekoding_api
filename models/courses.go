package models

import (
	"time"

	"gorm.io/gorm"
)

// Course is type for implemnting course model.
type Course struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:50;unique" json:"title"`
	Logo        string         `gorm:"size:255;null" json:"logo"`
	Description string         `gorm:"size:255" json:"description"`
	Content     string         `gorm:"type:text" json:"content"`
	IsActive    bool           `gorm:"default:false" json:"is_active"`
	OpenedAt    time.Time      `gorm:"size:10" json:"opened_at"`
	ClosedAt    time.Time      `gorm:"size:10" json:"closed_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Mentors     []*User        `gorm:"many2many:courses_users_mentor" json:"mentors"`
	Members     []*User        `gorm:"many2many:courses_users_member" json:"members"`
}

// CreateNewCourse is function for create new course.
func CreateNewCourse(course *Course) error {
	return DB().Create(course).Error
}

// GetCourseByID is function to get course by id.
func GetCourseByID(id int, isActive bool) (Course, error) {
	var course Course
	err := DB().Model(&Course{}).Where("id = ? AND is_active = ?", id, isActive).
		Preload("Mentors").Preload("Members").First(&course).Error
	return course, err
}

// GetCourseByIDNotParam
func GetCourseByIDNotParam(id int) (Course, error) {
	var course Course
	err := DB().Model(&Course{}).Where("id = ?", id).
		Preload("Mentors").Preload("Members").First(&course).Error
	return course, err
}

// GetAllCourse is function to get all course.
func GetAllCourse(isActive bool) []Course {
	var courses []Course
	DB().Model(&Course{}).Where("is_active = ?", isActive).
		Preload("Mentors").Preload("Members").Find(&courses)

	return courses
}
