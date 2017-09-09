package server

import (
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wen-bing/go-enterprise-web-toolkit/core"
	"github.com/wen-bing/go-enterprise-web-toolkit/core/db"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user"
	"log"
	"os"
)

type ServerConfig struct {
	Port int         `json:"port"`
	DB   db.DBConfig `json:"db"`
}

type ApplicationServer struct {
	modules           []core.Module
	config            ServerConfig
	router            *gin.Engine
	db                *gorm.DB
	jwtAuthMiddleware *jwt.GinJWTMiddleware
}

func New(config ServerConfig) *ApplicationServer {
	s := ApplicationServer{}
	s.config = config
	s.setupDatabase(s.config.DB)

	//load modules
	s.loadModules()

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

func (s *ApplicationServer) setupDatabase(dbConfig db.DBConfig) {
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
	s.jwtAuthMiddleware = userModule.GetJWTAuthMiddleware()
	s.modules = append(s.modules, userModule)
}

func (s *ApplicationServer) setupModuleRoutersV1() {
	//public router
	v1 := s.router.Group("api/v1")
	{
		//setup each modules router for v1
		for _, m := range s.modules {
			m.SetupRouterV1(v1)
		}
	}

	//security routers
	securityV1 := s.router.Group("api/s/v1")
	securityV1.Use(s.jwtAuthMiddleware.MiddlewareFunc())
	{
		//setup each modules router for v1
		for _, m := range s.modules {
			m.SetupSecurityRouterV1(securityV1)
		}
	}

}

func (s ApplicationServer) Start() {
	defer s.db.Close()
	s.router.Run(":8080")
}
