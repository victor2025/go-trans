package services

import (
	"encoding/json"
	"fmt"
	"go-trans/http/response"
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/dto"
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gorm.io/gorm"
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
	deviceName := GetServiceContext().ConfigService.GetOrDefault(consts.SelfDeviceName, deviceInfo.DeviceId)
	deviceInfo.DeviceName = deviceName
	return deviceInfo
}

// BePaired 以接收者的身份配对设备
func (s *DeviceService) BePaired(deviceInfo *entity.DeviceInfo, pairCode string) error {
	selfDeviceInfo := s.GetSelfDeviceInfo()
	if selfDeviceInfo.PairCode != strings.ToLower(pairCode) {
		return fmt.Errorf("invalid pair code")
	}
	strings.EqualFold(selfDeviceInfo.DeviceId, deviceInfo.DeviceId)
	var device *entity.DeviceInfo
	device = s.GetConnectedDeviceById(deviceInfo.DeviceId)
	if device == nil {
		device = entity.GetNewReceiveDeviceInfo(deviceInfo.DeviceId, deviceInfo.DeviceName)
	}
	device.Connected = true
	err := s.db.Save(device).Error
	utils.HandleError(err, utils.PanicOnError)
	err = s.RefreshSelfPairCode()
	utils.HandleError(err, utils.PanicOnError)
	return nil
}

func (s *DeviceService) Pair(deviceScanInfo *dto.DeviceScanInfo, pairCode string) error {
	ip := deviceScanInfo.Ip
	port := deviceScanInfo.Port
	deviceId := deviceScanInfo.DeviceId
	// 构建请求
	urlStr := fmt.Sprintf("http://%s:%s/device/inner/pair", ip, port)
	selfDeviceInfo := s.GetSelfDeviceInfo()
	// 补充端口号
	selfDeviceInfo.Port = GetServiceContext().ConfigService.GetOrDefault(consts.HttpServerPort, consts.DefaultHttpPort)
	selfDeviceInfo.TransmitPort = GetServiceContext().ConfigService.GetOrDefault(consts.TransServerPort, consts.DefaultHttpPort)
	deviceInfoBytes, err := json.Marshal(selfDeviceInfo)
	utils.HandleError(err)

	formData := url.Values{}
	formData.Set("deviceInfo", string(deviceInfoBytes))
	formData.Set("pairCode", pairCode)

	request, err := http.NewRequest("POST", urlStr, strings.NewReader(formData.Encode()))
	// 构建form数据
	// 发送请求
	utils.HandleError(err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("source", consts.DeviceScanIdentityParam)

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send pair request: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	utils.HandleError(err)
	respEntity := &response.Entity{}
	err = json.Unmarshal(body, respEntity)
	utils.HandleError(err)

	if respEntity.Success != consts.YES {
		return fmt.Errorf("pair failed: %v", respEntity.ErrMsg)
	}

	// 保存设备信息
	device := s.GetConnectedDeviceById(deviceId)
	if device == nil {
		device = entity.GetNewSendDeviceInfo(deviceId, "")
	}
	device.Connected = true
	err = s.db.Save(device).Error
	utils.HandleError(err)
	return nil
}

// RefreshSelfPairCode 刷新配对码
func (s *DeviceService) RefreshSelfPairCode() error {
	deviceInfo := s.GetSelfDeviceInfo()
	deviceInfo.RefreshPairCode()
	return s.db.Save(&deviceInfo).Error
}
