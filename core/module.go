package core

import (
	"github.com/gin-gonic/gin"
)

/**
Module defenition:
setup router
setup model
db migration

*/
type Module interface {
	/**
	Setup router for this module
	*/
	SetupRouter(r *gin.RouterGroup)
}
