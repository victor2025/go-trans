package system

import "github.com/gin-gonic/gin"

/**
  @author: victor2022
  @since: 2025/1/26
*/

func RegisterRouter(router *gin.RouterGroup) {
	// transmit server
	router.POST("/startTransmitServer", startTransmitServer)
	router.POST("/stopTransmitServer", stopTransmitServer)
	// settings
	router.POST("/updateSetting", updateSetting)
	router.GET("/getSetting", getSetting)
}
