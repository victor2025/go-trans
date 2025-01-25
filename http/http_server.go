package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-trans/http/api/system"
	"go-trans/http/api/task"
	"go-trans/services"
	"go-trans/utils"
)

/*
*

	@author: victor2022
	@since: 2025/1/5
*/
func StartHttpServer() {
	engine := initGinEngine()
	serverPort := services.GetServiceContext().ConfigService.GetOrDefault("http.server.port", "8080")
	// 读取配置
	err := engine.Run(":" + serverPort)
	utils.HandleError(err, utils.ExitOnErr)
}

func initGinEngine() *gin.Engine {
	engine := gin.Default()
	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// system
	systemApi := engine.Group("/system")
	system.RegisterRouter(systemApi)

	// task
	taskApi := engine.Group("/task")
	task.RegisterRouter(taskApi)
	return engine
}
