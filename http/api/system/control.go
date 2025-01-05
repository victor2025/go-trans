package system

import (
	"github.com/gin-gonic/gin"
	"go-trans/context"
	"go-trans/http/response"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

func RegisterRouter(router *gin.RouterGroup) {
	router.POST("/startTransmitServer", startTransmitServer)
	router.POST("/stopTransmitServer", stopTransmitServer)
}

func startTransmitServer(c *gin.Context) {
	context.GetSystemContext().StartReceiveServer()
	response.NewSuccessResponse(c, "start success")
}

func stopTransmitServer(c *gin.Context) {
	context.GetSystemContext().StopReceiveServer()
	response.NewSuccessResponse(c, "stop success")
}
