package services

import (
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/dto"
	"go-trans/pkg/models/entity"
	transHandler "go-trans/pkg/transmit/handlers"
	"go-trans/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/13
*/
type SendTaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *SendTaskService {
	taskService := &SendTaskService{
		db: db,
	}

	return taskService
}

// CreateSendTask 创建发送任务
func (s *SendTaskService) CreateSendTask(path, receiverId string) error {
	taskInfo, err := entity.GetNewSendTaskInfo(path, receiverId)
	if err != nil {
		return err
	}
	tx := s.db.Create(&taskInfo)
	if tx.Error != nil {
		return tx.Error
	}
	log.Printf("create send task success, path: %s, receiverId: %s", path, receiverId)
	return nil
}

// UpdateSendTask 更新任务
func (s *SendTaskService) UpdateSendTask(info *entity.SendTaskInfo) error {
	tx := s.db.Save(info)
	if tx.Error != nil {
		return tx.Error
	}
	log.Printf("update send task success, taskId: %s, status: %v", info.TaskId, info.Status)
	return nil
}

// GetSendTasks 按页查询任务列表
func (s *SendTaskService) GetSendTasks(page, size int, order string) ([]*entity.SendTaskInfo, error) {
	if order == "" {
		order = "id ASC"
	}
	offset := (page - 1) * size
	var tasks []*entity.SendTaskInfo
	tx := s.db.Order(order).Offset(offset).Limit(size).Find(&tasks)
	return tasks, tx.Error
}

// GetSendTaskDtos 取发送任务
func (s *SendTaskService) GetSendTaskDtos(page, size int, order string) ([]*dto.SendTaskDto, error) {
	if order == "" {
		order = "id ASC"
	}
	offset := (page - 1) * size
	var tasks []*entity.SendTaskInfo
	tx := s.db.Where("status = ?", consts.Waiting).Order(order).Offset(offset).Limit(size).Find(&tasks)
	utils.HandleError(tx.Error, utils.PanicOnError)
	result := make([]*dto.SendTaskDto, 0)
	for _, task := range tasks {
		var deviceInfo *entity.DeviceInfo
		tx := s.db.Find(&deviceInfo, "device_id = ? AND connected = ?", task.ReceiverId, true)
		if tx.RowsAffected == 0 || deviceInfo == nil {
			task.Status = consts.Fail
			s.UpdateSendTask(task)
			continue
		}
		result = append(result, &dto.SendTaskDto{
			Task:   task,
			Device: deviceInfo,
		})
	}
	return result, nil
}

// SendTaskProcessor 发送任务处理器
type SendTaskProcessor struct {
	isOn            bool
	scanFunc        func() ([]*dto.SendTaskDto, error)
	callbackFunc    func(*dto.SendTaskDto)
	sendTaskService *SendTaskService
	busService      *BusService
}

func NewSendTaskProcessor(sendTaskService *SendTaskService, busService *BusService) SendTaskProcessor {
	return SendTaskProcessor{
		sendTaskService: sendTaskService,
		busService:      busService,
	}
}

func (s *SendTaskProcessor) Start() {
	// 创建一个时间间隔为 200ms 的 ticker
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop() // 确保在退出时停止 ticker
	s.isOn = true
	for range ticker.C {
		if !s.isOn {
			break
		}
		tasks, err := s.scanSendTasks()
		utils.HandleError(err)
		if tasks == nil || len(tasks) == 0 {
			continue
		}
		for _, task := range tasks {
			log.Printf("start to process send task %+v", task)
			task.Task.Status = consts.Processing
			err := s.sendTaskService.UpdateSendTask(task.Task)
			utils.HandleError(err)
			s.startSendTask(task)
		}
	}
	log.Println("stop send task processor by signal")
}

func (s *SendTaskProcessor) Stop() {
	s.isOn = false
}

func (s *SendTaskProcessor) scanSendTasks() ([]*dto.SendTaskDto, error) {
	return s.sendTaskService.GetSendTaskDtos(1, 10, "")
}

func (s *SendTaskProcessor) startSendTask(task *dto.SendTaskDto) {
	handler := transHandler.NewSendHandler(task, func(taskDto *dto.SendTaskDto) {
		err := serviceContext.BusService.PostMsg(taskDto)
		utils.HandleError(err)
	})
	go handler.Handle()
}
