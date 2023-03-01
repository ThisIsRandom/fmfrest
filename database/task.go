package database

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       *string     `gorm:"not null" json:"title"`
	Description *string     `gorm:"not null" json:"description"`
	Images      []TaskImage `json:"images"`
	UserID      *int        `gorm:"not null" json:"userId"`
	User        User        `json:"user"`
}

type TaskImage struct {
	gorm.Model
	Uri string `json:"uri"`
}
