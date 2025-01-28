package system

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/services"
)

/**
@author: victor2022
@since: 2025/1/5
*/

// 启动文件接收服务器
func startTransmitServer(c *gin.Context) {
	services.GetServiceContext().StartReceiveServer()
	response.NewSuccessResponse(c, "start success")
}

// 关闭文件接收服务器
func stopTransmitServer(c *gin.Context) {
	services.GetServiceContext().StopReceiveServer()
	response.NewSuccessResponse(c, "stop success")
}
