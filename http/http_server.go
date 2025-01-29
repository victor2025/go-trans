package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-trans/http/api/device"
	"go-trans/http/api/system"
	"go-trans/http/api/task"
	"go-trans/pkg/models/consts"
	"go-trans/services"
	"go-trans/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

/**
@author: victor2022
@since: 2025/1/5
*/

// http服务器
var httpServer *http.Server

// StartHttpServer 启动http服务器
func StartHttpServer() {
	engine := initGinEngine()
	portStr := services.GetServiceContext().ConfigService.GetOrDefault(consts.HttpServerPort, "9210")
	port, _ := strconv.Atoi(portStr)
	httpServer = &http.Server{
		Handler: engine,
	}
	go func() {
		var err error
		for retryCnt := 0; retryCnt < 10; retryCnt++ {
			// 启动服务器
			log.Printf("try to start http server at port: %d\n", port)
			httpServer.Addr = fmt.Sprintf(":%d", port)
			err := httpServer.ListenAndServe()
			if err == nil {
				break
			}
			if errors.Is(err, http.ErrServerClosed) {
				break
			}
			utils.HandleError(err)
			port += 1
		}
		utils.HandleError(err, utils.PanicOnError)
	}()
}

// ShutdownHttpServer 关闭http服务器
func ShutdownHttpServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := httpServer.Shutdown(ctx)
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
