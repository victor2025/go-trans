package services

import (
	"go-trans/pkg/models/dto"
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
	"log"
)

/**
  @author: victor2022
  @since: 2025/1/13
*/

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

// CreateSendTask 创建发送任务
func (s *TaskService) CreateSendTask(path, receiverId string) error {
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
func (s *TaskService) UpdateSendTask(info *entity.SendTaskInfo) error {
	tx := s.db.Save(info)
	if tx.Error != nil {
		return tx.Error
	}
	log.Printf("update send task success, taskId: %s", info.TaskId)
	return nil
}

// GetSendTasks 按页查询任务列表
func (s *TaskService) GetSendTasks(page, size int, order string) ([]*entity.SendTaskInfo, error) {
	if order == "" {
		order = "id ASC"
	}
	offset := (page - 1) * size
	var tasks []*entity.SendTaskInfo
	tx := s.db.Order(order).Offset(offset).Limit(size).Find(&tasks)
	return tasks, tx.Error
}

// GetSendTaskDtos 取发送任务
func (s *TaskService) GetSendTaskDtos(page, size int, order string) ([]*dto.SendTaskDto, error) {
	tasks, err := s.GetSendTasks(page, size, order)
	utils.HandleError(err, utils.PanicOnError)
	result := make([]*dto.SendTaskDto, len(tasks))
	for _, task := range tasks {
		var deviceInfo *entity.DeviceInfo
		s.db.Find(&deviceInfo, "device_id = ?", task.ReceiverId)
		if deviceInfo == nil {
			continue
		}
		result = append(result, &dto.SendTaskDto{
			Task:   task,
			Device: deviceInfo,
		})
	}
	return result, nil
}
