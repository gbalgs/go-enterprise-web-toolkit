package core

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/**
Module defenition:
setup router
setup model
db migration

*/
type Module interface {
	SetupRouter(r *gin.RouterGroup)
	//
	MigrateSchema(db *gorm.DB)
	//UpMigration()
	//DownMigration()
}
