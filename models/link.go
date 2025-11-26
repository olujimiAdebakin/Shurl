package models

import "gorm.io/gorm"


type Link struct{

	// @dev gorm.Model is embedded to provide standard ID, CreatedAt, UpdatedAt, and DeletedAt fields.
	gorm.Model

	ShortCode string `gorm:"uniqueIndex"`
	OriginalURL string `gorm:"NOT NULL"`
	Hash string `gorm:"uniqueIndex;NOT NULL"`
	Clicks int `gorm:"default:0"`
	Favicon *string 
	UserID uint `gorm:"default:0"`
}