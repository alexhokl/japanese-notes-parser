package database

import (
	"gorm.io/gorm"
)

type StringArray []string

type Entry struct {
	gorm.Model
	Kanji   string
	Kana    string
	English StringArray `gorm:"type:jsonb"`
	Labels  StringArray `gorm:"type:jsonb"`
}

