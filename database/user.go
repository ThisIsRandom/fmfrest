package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                *string             `gorm:"type:varchar(255);unique;not null;" json:"email"`
	Password             *string             `json:"password" gorm:"not null"`
	ContactInformationID *int                `json:"contactInformationId"`
	ContactInformation   *ContactInformation `json:"contactInformation"`
	ProfileID            *int                `json:"profileId"`
	Profile              Profile             `json:"profile"`
	Advertisements       []Advertisement     `json:"advertisements"`
}

type ContactInformation struct {
	gorm.Model
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Address string `json:"address"`
	Postal  string `json:"postal"`
}

type Profile struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Avatar      string  `json:"avatar"`
	Images      []Image `json:"images"`
	RoleID      *int    `json:"roleId"`
	Role        Role    `json:"role"`
}

type Image struct {
	gorm.Model
	Description *string `json:"description" gorm:"not null"`
	Uri         *string `json:"uri"`
	ProfileID   *int    `json:"profileId"`
}

type Role struct {
	gorm.Model
	Name *string `gorm:"not null;unique" json:"name"`
}
