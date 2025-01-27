package services

import (
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
)

/*
*

	@author: victor2022
	@since: 2025/1/27
*/
type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService(db *gorm.DB) *DeviceService {
	return &DeviceService{
		db: db,
	}
}

func (s *DeviceService) GetConnectedDeviceById(deviceId string) *entity.DeviceInfo {
	var deviceInfo *entity.DeviceInfo
	tx := s.db.Find(&deviceInfo, "device_id = ? AND connected = ?", deviceId, true)
	utils.HandleError(tx.Error)
	if tx.RowsAffected == 0 {
		return nil
	}
	return deviceInfo
}

func (s *DeviceService) GetSelfDeviceInfo() *entity.DeviceInfo {
	var deviceInfo *entity.DeviceInfo
	tx := s.db.Find(&deviceInfo, "device_name = ?", consts.SelfDeviceName)
	if tx.RowsAffected == 0 {
		deviceInfo = s.initSelfDeviceInfo()
	}
	return deviceInfo
}

func (s *DeviceService) initSelfDeviceInfo() *entity.DeviceInfo {
	deviceInfo := entity.GetNewDeviceInfo("localhost", "")
	deviceInfo.DeviceName = consts.SelfDeviceName
	s.db.Create(&deviceInfo)
	return deviceInfo
}
