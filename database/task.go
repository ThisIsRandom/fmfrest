package database

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title          *string         `gorm:"not null" json:"title"`
	Description    *string         `gorm:"not null" json:"description"`
	Images         []TaskImage     `json:"images"`
	UserID         *int            `gorm:"not null" json:"userId"`
	User           User            `json:"user"`
	MessageStreams []MessageStream `json:"messageStream"`
}

type MessageStream struct {
	gorm.Model
	Messages []Message
	TaskID   int `json:"taskId"`
}

type Message struct {
	gorm.Model
	Text            *string `gorm:"not null"`
	MessageStreamID int     `json:"messageStreamId"`
}

type TaskImage struct {
	gorm.Model
	Uri    string `json:"uri"`
	TaskID int    `json:"taskId"`
}
