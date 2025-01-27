package services

import (
	"fmt"
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
	"strings"
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
	tx := s.db.Find(&deviceInfo, "mode = ?", consts.SelfMode)
	if tx.RowsAffected == 0 {
		deviceInfo = s.initSelfDeviceInfo()
	}
	return deviceInfo
}

func (s *DeviceService) PairDeviceForReceive(deviceId, deviceName, pairCode string) error {
	selfDeviceInfo := s.GetSelfDeviceInfo()
	if selfDeviceInfo.PairCode != strings.ToLower(pairCode) {
		return fmt.Errorf("invalid pair code")
	}
	strings.EqualFold(selfDeviceInfo.DeviceId, deviceId)
	var device *entity.DeviceInfo
	device = s.GetConnectedDeviceById(deviceId)
	if device == nil {
		device = entity.GetNewReceiveDeviceInfo(deviceId, deviceName)
	}
	device.PairCode = pairCode
	device.Connected = true
	device.DeviceName = deviceName
	err := s.db.Save(device).Error
	utils.HandleError(err, utils.PanicOnError)
	err = s.RefreshSelfPairCode()
	utils.HandleError(err, utils.PanicOnError)
	return nil
}

func (s *DeviceService) initSelfDeviceInfo() *entity.DeviceInfo {
	deviceInfo := entity.GetNewDeviceInfo("localhost", "", consts.SelfMode)
	s.db.Create(&deviceInfo)
	return deviceInfo
}

func (s *DeviceService) RefreshSelfPairCode() error {
	deviceInfo := s.GetSelfDeviceInfo()
	deviceInfo.RefreshPairCode()
	return s.db.Save(&deviceInfo).Error
}
