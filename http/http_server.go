package http

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/api/system"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/
func StartHttpServer() {
	engine := initGinEngine()
	// todo 读取配置
	engine.Run(":8080")
}

func initGinEngine() *gin.Engine {
	engine := gin.Default()
	// system
	systemApi := engine.Group("/system")
	system.RegisterRouter(systemApi)

	return engine
}
