package services

import (
	"fmt"
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/dto"
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
	"strings"
)

/**
@author: victor2022
@since: 2025/1/27
*/

// DeviceService 设备服务
type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService(db *gorm.DB) *DeviceService {
	return &DeviceService{
		db: db,
	}
}

// GetConnectedDeviceById 根据id查询已连接设备
func (s *DeviceService) GetConnectedDeviceById(deviceId string) *entity.DeviceInfo {
	var deviceInfo *entity.DeviceInfo
	tx := s.db.Find(&deviceInfo, "device_id = ? AND connected = ?", deviceId, true)
	utils.HandleError(tx.Error)
	if tx.RowsAffected == 0 {
		return nil
	}
	return deviceInfo
}

// GetSelfDeviceInfo 查询自身设备信息
func (s *DeviceService) GetSelfDeviceInfo() *entity.DeviceInfo {
	var deviceInfo *entity.DeviceInfo
	tx := s.db.Find(&deviceInfo, "mode = ?", consts.SelfMode)
	if tx.RowsAffected == 0 {
		// 创建新的设备信息
		deviceInfo = entity.GetNewDeviceInfo("localhost", "", consts.SelfMode)
		s.db.Create(&deviceInfo)
	}
	return deviceInfo
}

// PairDeviceForReceive 以接收者的身份配对设备
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

// RefreshSelfPairCode 刷新配对码
func (s *DeviceService) RefreshSelfPairCode() error {
	deviceInfo := s.GetSelfDeviceInfo()
	deviceInfo.RefreshPairCode()
	return s.db.Save(&deviceInfo).Error
}

// DeviceScanProcessor 设备扫描处理器
type DeviceScanner struct {
	scanResult map[string]dto.DeviceScanInfo
}

func NewDeviceScanner() *DeviceScanner {
	processor := &DeviceScanner{
		scanResult: make(map[string]dto.DeviceScanInfo),
	}
	return processor
}
