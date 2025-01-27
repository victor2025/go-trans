package system

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/services"
)

/*
*

	@author: victor2022
	@since: 2025/1/5
*/
func startTransmitServer(c *gin.Context) {
	services.GetServiceContext().StartReceiveServer()
	response.NewSuccessResponse(c, "start success")
}

func stopTransmitServer(c *gin.Context) {
	services.GetServiceContext().StopReceiveServer()
	response.NewSuccessResponse(c, "stop success")
}
