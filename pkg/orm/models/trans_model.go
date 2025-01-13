package models

import (
	"errors"
	"github.com/google/uuid"
	"go-trans/utils"
	"path/filepath"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/13
*/

type FileType string
type TaskStatus int

const (
	file       FileType   = "file"
	dir        FileType   = "dir"
	Waiting    TaskStatus = 0
	Processing TaskStatus = 1
	Finished   TaskStatus = 2
	Fail       TaskStatus = -1
)

type BaseTaskInfo struct {
	BaseModel
	FileName string   `gorm:"not null" json:"file_name"`
	FilePath string   `gorm:"not null" json:"file_path"`
	FileType FileType `json:"type"`
	TaskId   string   `gorm:"unique; not null" json:"task_id"`
	ErrorMsg string   `json:"error_msg"`
}

// SendTaskInfo 发送任务信息
type SendTaskInfo struct {
	BaseTaskInfo
	Progress   float32    `gorm:"index:idx_send_progress_status" json:"progress"`
	Status     TaskStatus `gorm:"index:idx_send_progress_status" json:"status"`
	ReceiverId string     `gorm:"index:idx_receiver" json:"receiver_id"`
}

// ReceiveTaskInfo 接收任务信息
type ReceiveTaskInfo struct {
	BaseTaskInfo
	Progress float32    `gorm:"index:idx_receive_progress_status" json:"progress"`
	Status   TaskStatus `gorm:"index:idx_receive_progress_status" json:"status"`
	SenderId string     `gorm:"index:idx_sender" json:"sender_id"`
}

func GetNewSendTaskInfo(path, receiverId string) (*SendTaskInfo, error) {
	if !utils.Exists(path) {
		return nil, errors.New("file not exists")
	}
	absPath, err := filepath.Abs(path)
	utils.HandleError(err)
	_, name := filepath.Split(absPath)
	var fileType FileType
	if utils.IsDir(absPath) {
		fileType = dir
	} else {
		fileType = file
	}
	return &SendTaskInfo{
		BaseTaskInfo: BaseTaskInfo{
			BaseModel: BaseModel{
				GmtCreate: time.Now(),
				GmtModify: time.Now(),
			},
			FileName: name,
			FilePath: absPath,
			FileType: fileType,
			TaskId:   uuid.New().String(),
		},
		ReceiverId: receiverId,
		Status:     Waiting,
	}, nil
}
