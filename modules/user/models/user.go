package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name        string `gorm:"size:255";index`
	Password    string `gorm:"size:255"`
	Email       string `gorm:"size:255;unique_index"`
	MobilePhone string `gorm:"size:64;index"`
	Profile     UserProfile
	ProfileID   int
}
