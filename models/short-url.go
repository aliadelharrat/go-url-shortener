package models

import "gorm.io/gorm"

type ShortURL struct {
	gorm.Model
	URL    string
	SURL   string `gorm:"column:surl"`
	Visits int
}