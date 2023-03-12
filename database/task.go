package database

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title          *string         `gorm:"not null" json:"title"`
	Description    *string         `gorm:"not null" json:"description"`
	Images         []TaskImage     `json:"images"`
	UserID         *int            `gorm:"not null" json:"userId"`
	User           User            `json:"user"`
	MessageStreams []MessageStream `json:"messageStreams"`
}

type MessageStream struct {
	gorm.Model
	Messages []Message `json:"messages"`
	TaskID   int       `json:"taskId"`
	UserID   int       `json:"userId"`
	User     User      `json:"user"`
}

type Message struct {
	gorm.Model
	Text            *string `gorm:"not null" json:"text"`
	MessageStreamID *int    `json:"messageStreamId"`
	UserID          *int    `json:"userId" gorm:"not null"`
	User            User    `json:"user"`
}

type TaskImage struct {
	gorm.Model
	Uri    string `json:"uri"`
	TaskID int    `json:"taskId"`
}
