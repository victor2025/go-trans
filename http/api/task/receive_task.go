package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/pkg/models/consts"
	"go-trans/services"
	"go-trans/utils"
	"strconv"
	"strings"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/

func pageReceiveTasks(c *gin.Context) {
	page := c.Param("page")
	size := c.Param("size")
	completed := c.Query("completed")
	statusListStr := ""
	if completed == consts.YES {
		statusListStr = fmt.Sprintf("%d,%d", consts.Finished, consts.Fail)
	} else if completed == consts.NO {
		statusListStr = fmt.Sprintf("%d,%d", consts.Waiting, consts.Processing)
	} else {
		statusListStr = fmt.Sprintf("%d,%d,%d,%d", consts.Waiting, consts.Processing, consts.Finished, consts.Fail)
	}
	statusList := strings.Split(statusListStr, ",")
	pageInt, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(size)

	tasks, err := services.GetServiceContext().ReceiveTaskService.GetReceiveTasks(pageInt, pageSize, statusList, "gmt_create DESC")
	utils.HandleError(err, func(args ...interface{}) {
		response.NewFailResponse(c, "", err.Error())
		panic(err.Error())
	})
	count, err := services.GetServiceContext().ReceiveTaskService.CountReceiveTasks(statusList)
	utils.HandleError(err, func(args ...interface{}) {
		response.NewFailResponse(c, "", err.Error())
		panic(err.Error())
	})
	response.NewSuccessResponse(c, map[string]interface{}{
		"tasks": tasks,
		"count": count,
	})
}

func deleteReceiveTask(c *gin.Context) {
	taskId := c.PostForm("taskId")
	err := services.GetServiceContext().ReceiveTaskService.DeleteReceiveTask(taskId)
	utils.HandleError(err, utils.PanicOnError)
	response.NewSuccessResponse(c, "success")
}
