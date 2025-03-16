package task

import "github.com/gin-gonic/gin"

/**
@author: victor2022
@since: 2025/1/19
*/

func RegisterRouter(router *gin.RouterGroup) {
	// send
	router.POST("/createSendTask", createSendTask)
	router.POST("/deleteSendTask", deleteSendTask)
	router.GET("/pageSendTasks/:page/:size", pageSendTasks)
	// receive
	router.GET("/pageReceiveTasks/:page/:size", pageReceiveTasks)
	router.POST("/deleteReceiveTask", deleteReceiveTask)

}
