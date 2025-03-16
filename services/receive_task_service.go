package services

import (
	"go-trans/pkg/models/entity"
	"gorm.io/gorm"
	"log"
)

/*
*

	@author: victor2022
	@since: 2025/3/16
*/
type ReceiveTaskService struct {
	db *gorm.DB
}

func NewReceiveTaskService(db *gorm.DB) *ReceiveTaskService {
	return &ReceiveTaskService{
		db: db,
	}
}

// UpdateReceiveTask 更新任务
func (s *ReceiveTaskService) UpdateReceiveTask(info *entity.ReceiveTaskInfo) error {
	tx := s.db.Save(info)
	if tx.Error != nil {
		return tx.Error
	}
	log.Printf("update receive task success, taskId: %s, status: %v", info.TaskId, info.Status)
	return nil
}

// GetReceiveTasks 按页查询任务列表
func (s *ReceiveTaskService) GetReceiveTasks(page, size int, taskStatusList []string, order string) ([]*entity.ReceiveTaskInfo, error) {
	if order == "" {
		order = "gmt_create ASC"
	}
	offset := (page - 1) * size
	var tasks []*entity.ReceiveTaskInfo
	tx := s.db.Order(order).Offset(offset).Limit(size).Find(&tasks, "status in ?", taskStatusList)
	return tasks, tx.Error
}

// CountReceiveTasks 统计满足条件的任务总数
func (s *ReceiveTaskService) CountReceiveTasks(taskStatusList []string) (int64, error) {
	var count int64
	tx := s.db.Model(&entity.SendTaskInfo{}).Where("status in ?", taskStatusList).Count(&count)
	return count, tx.Error
}
