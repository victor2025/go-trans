package app

import (
	"context"
	"github.com/mustafaturan/bus/v3"
	"go-trans/pkg/models/dto"
	"go-trans/services"
	"go-trans/utils"
	"log"
)

/**
  @author: victor2022
  @since: 2025/1/25
*/
func RegisterMsgConsumers() {
	registerSendTaskStatusMsgHandler()
}

// 发送任务状态监听器
func registerSendTaskStatusMsgHandler() {
	serviceContext := services.GetServiceContext()
	// 注册消息监听
	serviceContext.BusService.RegisterHandler("sendStatusMsgHandler", func(ctx context.Context, e bus.Event) {
		log.Printf("receive msg id: %v", e.ID)
		taskDto := e.Data.(*dto.SendTaskDto)
		err := serviceContext.SendTaskService.UpdateSendTask(taskDto.Task)
		utils.HandleError(err)
	})
}
