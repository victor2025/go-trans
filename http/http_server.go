package http

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/api/device"
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
	// system
	systemApi := engine.Group("/system")
	system.RegisterRouter(systemApi)

	// device
	deviceApi := engine.Group("/device")
	device.RegisterRouter(deviceApi)

	// task
	taskApi := engine.Group("/task")
	task.RegisterRouter(taskApi)
	return engine
}
