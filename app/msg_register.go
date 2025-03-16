package app

import (
	"context"
	"github.com/mustafaturan/bus/v3"
	"go-trans/pkg/models/dto"
	"go-trans/services"
	"go-trans/utils"
)

/*
*

	@author: victor2022
	@since: 2025/1/25
*/
func RegisterMsgConsumers() {
	registerSendTaskStatusMsgHandler()
}

// 发送任务状态监听器
func registerSendTaskStatusMsgHandler() {
	serviceContext := services.GetServiceContext()
	// 发送消息监听
	serviceContext.BusService.RegisterHandler("sendStatusMsgHandler", func(ctx context.Context, e bus.Event) {
		services.InfoF("receive msg id: %v", e.ID)
		taskDto, ok := e.Data.(*dto.SendTaskDto)
		if ok {
			err := serviceContext.SendTaskService.UpdateSendTask(taskDto.Task)
			utils.HandleError(err)
		}
	})
	// 接收消息监听
	serviceContext.BusService.RegisterHandler("receiveStatusMsgHandler", func(ctx context.Context, e bus.Event) {
		services.InfoF("receive msg id: %v", e.ID)
		taskDto, ok := e.Data.(*dto.ReceiveTaskDto)
		if ok {
			err := serviceContext.ReceiveTaskService.UpdateReceiveTask(taskDto.Task)
			utils.HandleError(err)
		}
	})
}
