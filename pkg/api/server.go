package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port        string
	Host        string
	Controllers []Controller
}

func StartWebServer(config *ServerConfig) error {
	router := gin.Default()
	router.Use(gin.ErrorLogger()) // TODO: decide what to actually send. Errors might leak sensitive info.

	err := router.SetTrustedProxies(nil)
	if err != nil {
		return err
	}

	registerControllers(router, config.Controllers)

	return router.Run(fmt.Sprintf("%s:%s", config.Host, config.Port))
}

type Controller interface {
	Register(router *gin.RouterGroup)
}

func registerControllers(router *gin.Engine, controllers []Controller) {
	group := router.Group("")
	for _, controller := range controllers {
		controller.Register(group)
	}
}
