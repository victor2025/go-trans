package task

import (
	"github.com/gin-gonic/gin"
	"go-trans/context"
	"go-trans/http/response"
	"go-trans/utils"
	"log"
	"strconv"
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

func pageSendTasks(c *gin.Context) {
	page := c.Param("page")
	size := c.Param("size")
	pageInt, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(size)
	tasks, err := context.GetServiceContext().TaskService.GetSendTasks(pageInt, pageSize, "id DESC")
	utils.HandleError(err, func() {
		response.NewFailResponse(c, "", err.Error())
		panic(err.Error())
	})
	response.NewSuccessResponse(c, tasks)
}
