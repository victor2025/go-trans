package app

import (
	"go-trans/pkg/runner"
	"go-trans/services"
	"log"
	"time"
)

/**
@author: victor2022
@since: 2025/1/25
*/

// 保存所有runner
var runners []runner.Runner

func RegisterProcessors() {
	registerSendTaskProcessor()
}

// 注册发送任务处理器
func registerSendTaskProcessor() {
	runners = make([]runner.Runner, 0)
	start := time.Now()
	sendTaskService := services.GetServiceContext().SendTaskService
	busService := services.GetServiceContext().BusService
	// 启动processor
	processor := services.NewSendTaskProcessor(sendTaskService, busService)
	go processor.Run()
	log.Printf("Send task processor started in %dms\n", time.Since(start).Milliseconds())
	runners = append(runners, processor)
}

// UnregisterProcessors 关闭所有处理器
func UnregisterProcessors() {
	for _, runningRunner := range runners {
		runningRunner.Stop()
	}
}
