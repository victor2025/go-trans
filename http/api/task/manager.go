package task

import "github.com/gin-gonic/gin"

/**
  @author: victor2022
  @since: 2025/1/19
*/
func RegisterRouter(router *gin.RouterGroup) {
	router.POST("/createSendTask", createSendTask)
}
