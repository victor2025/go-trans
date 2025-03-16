package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/pkg/models/consts"
	"go-trans/services"
	"go-trans/utils"
	"log"
	"strconv"
	"strings"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/

func createSendTask(c *gin.Context) {
	path := c.PostForm("path")
	receiver := c.PostForm("receiver")
	err := services.GetServiceContext().SendTaskService.CreateSendTask(path, receiver)
	utils.HandleError(err, func(args ...interface{}) {
		response.NewFailResponse(c, "", err.Error())
		log.Printf("createSendTask err, errMsg: %s\n", err.Error())
		panic(err.Error())
	})
	response.NewSuccessResponse(c, "success")
}

func pageSendTasks(c *gin.Context) {
	page := c.Param("page")
	size := c.Param("size")
	statusListStr := c.Query("statusList")
	if strings.TrimSpace(statusListStr) == "" {
		statusListStr = fmt.Sprintf("%d,%d", consts.Waiting, consts.Processing)
	}
	statusList := strings.Split(statusListStr, ",")
	pageInt, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(size)

	tasks, err := services.GetServiceContext().SendTaskService.GetSendTasks(pageInt, pageSize, statusList, "gmt_create DESC")
	utils.HandleError(err, func(args ...interface{}) {
		response.NewFailResponse(c, "", err.Error())
		panic(err.Error())
	})
	response.NewSuccessResponse(c, tasks)
}
