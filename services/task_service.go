package services

import (
	"go-trans/pkg/orm/models"
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
	taskInfo, err := models.GetNewSendTaskInfo(path, receiverId)
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
func (s *TaskService) UpdateSendTask(info *models.SendTaskInfo) error {
	tx := s.db.Save(info)
	if tx.Error != nil {
		return tx.Error
	}
	log.Printf("update send task success, taskId: %s", info.TaskId)
	return nil
}

// GetTasks 按页查询任务列表
func (s *TaskService) GetTasks(page, size int, order string) ([]*models.SendTaskInfo, error) {
	if order == "" {
		order = "id ASC"
	}
	offset := (page - 1) * size
	var tasks []*models.SendTaskInfo
	tx := s.db.Order(order).Offset(offset).Limit(size).Find(&tasks)
	return tasks, tx.Error
}
