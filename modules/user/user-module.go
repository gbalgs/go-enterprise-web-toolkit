package user

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/controllers"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/services"
	"time"
)

/**
Usage:
1. init inner component: models, services , controlers
2. setup router
3. expose public elements: models, services
*/

type UserModule struct {
	db                *gorm.DB
	jwtAuthMiddleware *jwt.GinJWTMiddleware
	userService       *services.UserService
	userController    *controllers.UserController
	//profileController controllers.ProfileController
}

func New(db *gorm.DB) *UserModule {
	m := UserModule{}
	m.db = db
	m.userService = services.NewUserService(db)
	m.userController = controllers.NewUserController(m.userService)

	//setup authentication
	m.jwtAuthMiddleware = &jwt.GinJWTMiddleware{
		Realm:         "gewt",
		Key:           []byte("gewt super security key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Authenticator: m.userController.JWTAuthenticator,
		Authorizator:  m.userController.JWTAuthorizator,
		PayloadFunc:   m.userController.JWTPlayloadFunc,
		Unauthorized:  m.userController.JWTUnauthorized,
	}
	return &m
}

func (m *UserModule) GetJWTAuthMiddleware() *jwt.GinJWTMiddleware {
	return m.jwtAuthMiddleware
}

func (m *UserModule) SetupRouterV1(r *gin.RouterGroup) {
	///users
	//need verify email or phone
	r.POST("/users/registration", m.userController.Registration)
	r.POST("/users/tokens", m.jwtAuthMiddleware.LoginHandler)
	r.DELETE("/users/tokens", m.userController.Logout)
}

func (m *UserModule) SetupSecurityRouterV1(r *gin.RouterGroup) {
	//for admin to create user
	//no need vrify email or phone
	r.POST("/users", m.userController.CreateUser)
	r.GET("/users", m.userController.GetUsers)
	r.GET("/users/:id", m.userController.GetUser)
	r.PUT("/users/:id", m.userController.EditUser)
	r.DELETE("/users/:id", m.userController.DeleteUser)

	//profile
	//r.GET("/users/:id/profile", m.profileController.GetUserProfile)
	//r.POST("/users/:id/profile", m.profileController.CreateUserProfile)
	//r.PUT("/users/:id/profile", m.profileController.EditUserProfile)
	//r.DELETE("/users/:id/profile", m.profileController.DeleteProfile)
}
