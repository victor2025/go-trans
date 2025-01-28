package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-trans/http/api/device"
	"go-trans/http/api/system"
	"go-trans/http/api/task"
	"go-trans/pkg/models/consts"
	"go-trans/services"
	"go-trans/utils"
	"log"
	"strconv"
)

/**
@author: victor2022
@since: 2025/1/5
*/

// StartHttpServer 启动http服务器
func StartHttpServer() {
	engine := initGinEngine()
	portStr := services.GetServiceContext().ConfigService.GetOrDefault(consts.HttpServerPort, "8080")
	port, _ := strconv.Atoi(portStr)
	var err error
	for retryCnt := 0; retryCnt < 10; retryCnt++ {
		// 启动服务器
		log.Printf("try to start http server at port: %d\n", port)
		err = engine.Run(fmt.Sprintf(":%d", port))
		if err == nil {
			break
		}
		port += 1
	}
	utils.HandleError(err, utils.PanicOnError)
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
