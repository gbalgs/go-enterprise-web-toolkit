package models

import (
	"github.com/wen-bing/go-enterprise-web-toolkit/core/db"
)

type User struct {
	db.Model
	Name        string
	Password    string
	Email       string
	MobilePhone string
	Profile     UserProfile
	ProfileID   int
}
