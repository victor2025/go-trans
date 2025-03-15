package http

import (
	"context"
	"errors"
	"fmt"
	"go-trans/http/api/device"
	"go-trans/http/api/system"
	"go-trans/http/api/task"
	"go-trans/pkg/models/consts"
	"go-trans/services"
	"go-trans/utils"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/**
@author: victor2022
@since: 2025/1/5
*/

// http服务器
var httpServer *http.Server

// StartHttpServer 启动http服务器
func StartHttpServer() int {
	engine := initGinEngine()
	portStr := services.GetServiceContext().ConfigService.GetOrDefault(consts.HttpServerPort, "9210")
	port, _ := strconv.Atoi(portStr)
	httpServer = &http.Server{
		Handler: engine,
	}
	// 创建channel用于获取最终启动的端口
	portChan := make(chan int, 1)
	go func() {
		var err error
		for retryCnt := 0; retryCnt < 10; retryCnt++ {
			// 启动服务器
			log.Printf("try to start http server at port: %d\n", port)
			httpServer.Addr = fmt.Sprintf(":%d", port)
			listener, err := net.Listen("tcp", httpServer.Addr)
			if err == nil {
				// 发送绑定的端口
				portChan <- port
				err = httpServer.Serve(listener)
				if errors.Is(err, http.ErrServerClosed) {
					break
				}
			}
			utils.HandleError(err)
			port += 1
		}
		utils.HandleError(err, utils.PanicOnError)
	}()
	// 返回最终启动的端口
	return <-portChan
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
