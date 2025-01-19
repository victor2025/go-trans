package dto

import (
	"go-trans/pkg/models/entity"
)

/**
  @author: victor2022
  @since: 2025/1/19
*/
type SendTaskDto struct {
	Task   *entity.SendTaskInfo
	Device *entity.DeviceInfo
}
