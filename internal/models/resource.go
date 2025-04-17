package models

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	MD5        string `gorm:"size:255;index;not null"`
	FileType   string `gorm:"not null"`
	FilePath   string `gorm:"not null"`
	StaticPath string `gorm:"not null"`
}
