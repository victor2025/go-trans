package http

import (
	"github.com/gin-gonic/gin"
	"go-trans/context"
	"go-trans/http/api/system"
	"go-trans/utils"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/
func StartHttpServer() {
	engine := initGinEngine()
	serverPort := context.GetServiceContext().ConfigService.GetOrDefault("http.server.port", "8080")
	// 读取配置
	err := engine.Run(":" + serverPort)
	utils.HandleError(err, utils.ExitOnErr)
}

func initGinEngine() *gin.Engine {
	engine := gin.Default()
	// system
	systemApi := engine.Group("/system")
	system.RegisterRouter(systemApi)

	return engine
}
