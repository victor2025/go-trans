package task

import (
	"go-trans/pkg/models/dto"
	transHandler "go-trans/pkg/transmit/handlers"
	"go-trans/utils"
	"log"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/
type SendTaskProcessor struct {
	isOn       bool
	scanFunc   func() (*[]dto.SendTaskDto, error)
	recordFunc func(*dto.SendTaskDto) error
}

func NewSendTaskProcessor(scanFunc func() (*[]dto.SendTaskDto, error), recordFunc func(info *dto.SendTaskDto) error) SendTaskProcessor {
	return SendTaskProcessor{
		scanFunc:   scanFunc,
		recordFunc: recordFunc,
	}
}

func (s *SendTaskProcessor) Start() {
	// 创建一个时间间隔为 200ms 的 ticker
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // 确保在退出时停止 ticker

	for range ticker.C {
		if !s.isOn {
			break
		}
		tasks, err := s.scanSendTasks()
		utils.HandleError(err)
		if tasks == nil || len(*tasks) == 0 {
			continue
		}
		for _, task := range *tasks {
			log.Printf("start to process send task %+v", task)
			s.startSendTask(&task)
		}
	}
	log.Println("stop send task processor by signal")
}

func (s *SendTaskProcessor) Stop() {
	s.isOn = false
}

func (s *SendTaskProcessor) scanSendTasks() (*[]dto.SendTaskDto, error) {
	return s.scanFunc()
}

func (s *SendTaskProcessor) startSendTask(task *dto.SendTaskDto) {
	path := task.Task.FilePath
	address := task.Device.Address
	port := task.Device.Port
	transHandler.NewSendHandler(address, port, path)
}

func (s *SendTaskProcessor) recordSendTask(task *dto.SendTaskDto) error {
	return s.recordFunc(task)
}
