package app

import (
	"go-trans/services"
	"log"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/25
*/
func RegisterProcessors() {
	registerSendTaskProcessor()
}

// 注册发送任务处理器
func registerSendTaskProcessor() {
	start := time.Now()
	sendTaskService := services.GetServiceContext().SendTaskService
	busService := services.GetServiceContext().BusService
	// 启动processor
	processor := services.NewSendTaskProcessor(sendTaskService, busService)
	go processor.Start()
	log.Printf("Send task processor started in %dms\n", time.Since(start).Milliseconds())
}
