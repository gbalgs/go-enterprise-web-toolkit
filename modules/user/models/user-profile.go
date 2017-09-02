package models

import "github.com/jinzhu/gorm"

type UserProfile struct {
	gorm.Model
	Company  string `gorm:"size:255"`
	Address1 string `gorm:size:1024`
	Address2 string `gorm:size:1024`
	Sex      string `gorm:size:16`
}
