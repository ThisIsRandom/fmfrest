package database

import "gorm.io/gorm"

type Advertisement struct {
	gorm.Model
	Image       string  `json:"image"`
	Activated   bool    `gorm:"not null;default:false;" json:"activated"`
	Description *string `gorm:"not null;" json:"description"`
	UserID      int     `json:"userId"`
}
