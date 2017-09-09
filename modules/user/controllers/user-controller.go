package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/models"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user/services"
	"net/http"
)

type UserController struct {
	s *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	c := UserController{}
	c.s = s
	return &c
}

/**
gin-jwt Callback function that should perform the authentication of the user based on userID and
password. Must return true on success, false on failure. Required.
Option return user id, if so, user id will be stored in Claim Array.
*/
func (u *UserController) JWTAuthenticator(userId string, password string, c *gin.Context) (string, bool) {
	id, err := u.s.Login(userId, password)
	if err != nil {
		return "", false
	} else {
		return id, true
	}
}

/**
gin-jwt Callback function that should perform the authorization of the authenticated user. Called
only after an authentication success. Must return true on success, false on failure.
Optional, default to success.
*/
func (u *UserController) JWTAuthorizator(userID string, c *gin.Context) bool {
	//TODO
	//verify user's permission
	if userID != "" {
		return true
	} else {
		return false
	}
}

// gin-jwt Callback function that will be called during login.
// Using this function it is possible to add additional payload data to the webtoken.
// The data is then made available during requests via c.Get("JWT_PAYLOAD").
// Note that the payload is not encrypted.
// The attributes mentioned on jwt.io can't be used as keys for the map.
// Optional, by default no additional data will be set.
func (m *UserController) JWTPlayloadFunc(userID string) map[string]interface{} {
	result := make(map[string]interface{})
	result["iss"] = "gewt.com"
	return result
}

/**
gin-jwt callback to handle handle unauthorized
*/
func (m *UserController) JWTUnauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (this *UserController) Registration(c *gin.Context) {
	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userObj := this.s.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{"user": userObj})
}

func (u *UserController) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, "email or password not sepcified")
		return
	}

	userObj, err := u.s.Login(email, password)
	if err != nil {
		c.JSON(http.StatusForbidden, "email or password mismatch")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userObj})
}

func (u *UserController) Logout(c *gin.Context) {

}

func (u *UserController) CreateUser(c *gin.Context) {

}

func (u *UserController) GetUsers(c *gin.Context) {

	c.JSON(200, gin.H{"data": "Hello GEWT"})
}

func (u *UserController) GetUser(c *gin.Context) {
	id, exist := c.Get("userID")
	if exist {
		obj, err := u.s.GetUserById(id.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch user by id failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": obj})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	}
}

func (u *UserController) EditUser(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Hello GEWT"})
}

func (u *UserController) DeleteUser(c *gin.Context) {

	c.JSON(200, gin.H{"data": "Hello GEWT"})
}
