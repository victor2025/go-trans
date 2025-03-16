package entity

import (
	"errors"
	"github.com/google/uuid"
	"go-trans/pkg/models/consts"
	"go-trans/utils"
	"path/filepath"
)

/**
  @author: victor2022
  @since: 2025/1/13
*/

type BaseTaskInfo struct {
	*BaseModel
	FileName string          `gorm:"not null" json:"fileName"`
	FilePath string          `gorm:"not null" json:"filePath"`
	FileType consts.FileType `json:"type"`
	Md5      string          `json:"md5"`
	TaskId   string          `gorm:"unique; not null" json:"taskId"`
	ErrorMsg string          `json:"error_msg"`
}

// SendTaskInfo 发送任务信息
type SendTaskInfo struct {
	*BaseTaskInfo
	Progress   float32           `gorm:"index:idx_send_progress_status" json:"progress"`
	Status     consts.TaskStatus `gorm:"index:idx_send_progress_status" json:"status"`
	ReceiverId string            `gorm:"index:idx_receiver" json:"receiverId"`
}

// ReceiveTaskInfo 接收任务信息
type ReceiveTaskInfo struct {
	*BaseTaskInfo
	Progress float32           `gorm:"index:idx_receive_progress_status" json:"progress"`
	Status   consts.TaskStatus `gorm:"index:idx_receive_progress_status" json:"status"`
	FileSize int64             `json:"fileSize"`
	SenderId string            `gorm:"index:idx_sender" json:"senderId"`
}

func GetNewSendTaskInfo(path, receiverId string) (*SendTaskInfo, error) {
	if !utils.Exists(path) {
		return nil, errors.New("file not exists")
	}
	absPath, err := filepath.Abs(path)
	utils.HandleError(err)
	_, name := filepath.Split(absPath)
	var fileType consts.FileType
	if utils.IsDir(absPath) {
		fileType = consts.Dir
	} else {
		fileType = consts.File
	}
	return &SendTaskInfo{
		BaseTaskInfo: &BaseTaskInfo{
			BaseModel: GetNewBaseModel(),
			FileName:  name,
			FilePath:  absPath,
			FileType:  fileType,
			TaskId:    uuid.New().String(),
		},
		ReceiverId: receiverId,
		Status:     consts.Waiting,
	}, nil
}

func GetNewReceiveTaskInfo() (*ReceiveTaskInfo, error) {
	return &ReceiveTaskInfo{
		BaseTaskInfo: &BaseTaskInfo{
			BaseModel: GetNewBaseModel(),
			FileType:  consts.File,
			TaskId:    uuid.New().String(),
		},
		Status: consts.Waiting,
	}, nil
}
