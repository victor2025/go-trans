package services

import (
	"fmt"
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/dto"
	"go-trans/pkg/models/entity"
	"go-trans/pkg/runner"
	transHandler "go-trans/pkg/transmit/handlers"
	"go-trans/utils"
	"gorm.io/gorm"
	"log"
	"runtime/debug"
	"time"
)

/*
*

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
	InfoF("update send task success, taskId: %s, status: %v", info.TaskId, info.Status)
	return nil
}

func (s *SendTaskService) DeleteSendTask(taskId string) error {
	tx := s.db.Delete(&entity.SendTaskInfo{}, "task_id=?", taskId)
	return tx.Error
}

// GetSendTasks 按页查询任务列表
func (s *SendTaskService) GetSendTasks(page, size int, taskStatusList []string, order string) ([]*entity.SendTaskInfo, error) {
	if order == "" {
		order = "gmt_create ASC"
	}
	offset := (page - 1) * size
	var tasks []*entity.SendTaskInfo
	tx := s.db.Order(order).Offset(offset).Limit(size).Find(&tasks, "status in ?", taskStatusList)
	return tasks, tx.Error
}

// CountSendTasks 统计满足条件的任务总数
func (s *SendTaskService) CountSendTasks(taskStatusList []string) (int64, error) {
	var count int64
	tx := s.db.Model(&entity.SendTaskInfo{}).Where("status in ?", taskStatusList).Count(&count)
	return count, tx.Error
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
		deviceInfo := GetServiceContext().DeviceService.GetConnectedDeviceById(task.ReceiverId)
		if deviceInfo == nil {
			task.Status = consts.Fail
			task.ErrorMsg = fmt.Sprintf("cannot find receiver, id:%v", task.ReceiverId)
			err := s.UpdateSendTask(task)
			utils.HandleError(err)
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
	runner.TickerRunner
	isOn            bool
	scanFunc        func() ([]*dto.SendTaskDto, error)
	callbackFunc    func(*dto.SendTaskDto)
	sendTaskService *SendTaskService
	busService      *BusService
}

func NewSendTaskProcessor(sendTaskService *SendTaskService, busService *BusService) *SendTaskProcessor {
	processor := &SendTaskProcessor{
		TickerRunner: runner.TickerRunner{
			Period: time.Millisecond * 200,
		},
		sendTaskService: sendTaskService,
		busService:      busService,
	}
	processor.Runner = processor
	return processor
}

func (p *SendTaskProcessor) MarkStarted() {
	p.isOn = true
}

func (p *SendTaskProcessor) IsOn() bool {
	return p.isOn
}

func (p *SendTaskProcessor) Handle() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("send task processor panic, err: %v, stack: %s", err, string(debug.Stack()))
		}
	}()

	tasks, err := p.scanSendTasks()
	utils.HandleError(err)
	if tasks == nil || len(tasks) == 0 {
		return
	}
	for _, task := range tasks {
		log.Printf("start to process send task %+v", task)
		task.Task.Status = consts.Processing
		err := p.sendTaskService.UpdateSendTask(task.Task)
		utils.HandleError(err)
		p.startSendTask(task)
	}
}

func (p *SendTaskProcessor) Stop() {
	p.isOn = false
}

func (p *SendTaskProcessor) scanSendTasks() ([]*dto.SendTaskDto, error) {
	return p.sendTaskService.GetSendTaskDtos(1, 10, "")
}

func (p *SendTaskProcessor) startSendTask(task *dto.SendTaskDto) {
	handler := transHandler.NewSendHandler(task, func(taskDto *dto.SendTaskDto) {
		err := serviceContext.BusService.PostMsg(taskDto)
		utils.HandleError(err)
	})
	go handler.Handle()
}
