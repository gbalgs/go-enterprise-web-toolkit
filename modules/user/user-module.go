package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/controllers"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/models"
)

/**
Usage:
1. init inner component: models, services , controlers
2. setup router
3. expose public elements: models, services
*/

type UserModule struct {
	controller controllers.UserController
}

func New(db *gorm.DB) UserModule {
	m := UserModule{}
	return m
}

func (m UserModule) MigrateSchema(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.UserProfile{})

	//setup relationship
	db.Model(&models.User{}).Related(&models.UserProfile{})
}

func (m UserModule) SetupRouter(r *gin.RouterGroup) {
	r.GET("/users", m.controller.GetAllUsers)
}
