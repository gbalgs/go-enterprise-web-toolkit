package user

import (
	"github.com/gin-gonic/gin"
)

/**
Usage:
1. init inner component: models, services , controlers
2. setup router
3. expose public elements: models, services
*/

type UserModule struct {
	controller UserController
}

func New() UserModule {
	m := UserModule{}
	return m
}
f
func (m UserModule) SetupRouter(r *gin.RouterGroup) {
	r.GET("/users", m.controller.GetAllUsers)
}
