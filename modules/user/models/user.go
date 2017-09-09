package models

import (
	"github.com/wen-bing/go-enterprise-web-toolkit/core/db"
)

type User struct {
	db.Model
	Name     string `json:"name" form:"name" binding:"exists,alphanum,min=4,max=255"`
	Password string `json:"password" form:"password" binding:"exists,min=8,max=255"`
	Email    string `json:"email" form:"email" binding:"email"`
	Phone    string `json:"phone" form:"phone"`
	Status   int    `json:"status" form "-"`
}

type UserBasicObject struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func UserModel2BasicObject(u User) UserBasicObject {
	o := UserBasicObject{}
	o.Id = u.Id
	o.Name = u.Name
	o.Email = u.Email
	o.Phone = u.Phone
	return o
}
