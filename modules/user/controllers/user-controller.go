package controllers

import "github.com/gin-gonic/gin"

type UserController struct {
}

func (u *UserController) GetAllUsers(c *gin.Context) {

	c.JSON(200, gin.H{"data": "Hello GEWT"})
}
