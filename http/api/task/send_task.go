package task

import (
	"github.com/gin-gonic/gin"
	"go-trans/context"
	"go-trans/http/response"
	"go-trans/utils"
	"log"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/

func createSendTask(c *gin.Context) {
	path := c.DefaultQuery("path", "")
	receiver := c.DefaultQuery("receiver", "")
	err := context.GetServiceContext().TaskService.CreateSendTask(path, receiver)
	utils.HandleError(err, func() {
		response.NewFailResponse(c, "", err.Error())
		log.Printf("createSendTask err, errMsg: %s\n", err.Error())
		panic(err.Error())
	})
	response.NewSuccessResponse(c, "success")
}
