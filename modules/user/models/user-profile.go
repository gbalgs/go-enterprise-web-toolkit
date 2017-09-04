package models

import "github.com/wen-bing/go-enterprise-web-toolkit/core/db"

type UserProfile struct {
	db.Model
	Company  string
	Address1 string
	Address2 string
	Sex      string
}
