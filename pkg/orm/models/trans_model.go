package models

/**
  @author: victor2022
  @since: 2025/1/13
*/

type BaseTaskInfo struct {
	BaseModel
	FileName string `gorm:"not null" json:"file_name"`
	FilePath string `gorm:"not null" json:"file_path"`
	TaskId   string `gorm:"unique; not null" json:"task_id"`
	ErrorMsg string `json:"error_msg"`
}

// SendTaskInfo 发送任务信息
type SendTaskInfo struct {
	BaseTaskInfo
	Progress   float32 `gorm:"index:idx_send_progress_status" json:"progress"`
	Status     int     `gorm:"index:idx_send_progress_status" json:"status"`
	ReceiverId string  `gorm:"index:idx_receiver" json:"receiver_id"`
}

// ReceiveTaskInfo 接收任务信息
type ReceiveTaskInfo struct {
	BaseTaskInfo
	Progress float32 `gorm:"index:idx_receive_progress_status" json:"progress"`
	Status   int     `gorm:"index:idx_receive_progress_status" json:"status"`
	SenderId string  `gorm:"index:idx_sender" json:"sender_id"`
}
