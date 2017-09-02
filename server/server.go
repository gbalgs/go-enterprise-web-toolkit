package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wen-bing/go-enterprise-web-toolkit/core"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user"
	"log"
	"os"
)

type ApplicationServer struct {
	modules []core.Module
	config  ServerConfig
	router  *gin.Engine
	db      *gorm.DB
}

func New(env string, configDir string) *ApplicationServer {
	s := ApplicationServer{}
	s.config = initConfig(env, configDir)

	s.setupDatabase(s.config.DB)

	//load modules
	s.loadModules()

	//setup modules' db schema

	//set modules' router
	s.setupRouters()
	return &s
}
func (s *ApplicationServer) setupRouters() {
	s.router = gin.New()
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.setupModuleRoutersV1()
	//setup static router
}

func (s *ApplicationServer) setupDatabase(dbConfig core.DBConfig) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name)
	var err error
	s.db, err = gorm.Open(dbConfig.Type, dbUrl)
	if err != nil {
		log.Printf("DB connection error: %v", err)
		os.Exit(-2)
	}
}

func (s *ApplicationServer) loadModules() {
	//user modules
	userModule := user.New(s.db)
	s.modules = append(s.modules, userModule)
}

func (s *ApplicationServer) setupModuleRoutersV1() {
	//public router
	v1 := s.router.Group("api/v1")
	{
		//setup each modules router for v1
		for _, m := range s.modules {
			m.SetupRouter(v1)
		}
	}
	//security routers
}

func (s ApplicationServer) Start() {
	s.router.Run(":8080")
}
