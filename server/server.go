package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wen-bing/go-enterprise-web-toolkit/core"
	"github.com/wen-bing/go-enterprise-web-toolkit/modules/user"
)

type ApplicationServer struct {
	modules []core.Module
	config  ServerConfig
	router  *gin.Engine
}

func New(env string) *ApplicationServer {
	s := ApplicationServer{}
	s.config = initConfig(env)
	s.initModules()
	s.initRouter()
	s.initModulesRouterV1()
	return &s
}

func (s *ApplicationServer) initModules() {
	//user modules
	userModule := user.New()
	s.modules = append(s.modules, userModule)
}

func (s *ApplicationServer) initRouter() {
	s.router = gin.New()
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
}

func (s *ApplicationServer) initModulesRouterV1() {
	v1 := s.router.Group("api/v1")
	{
		//setup each modules router for v1
		for _, m := range s.modules {
			m.SetupRouter(v1)
		}
	}
}

func (s ApplicationServer) Start() {
	s.router.Run(":8080")
}
