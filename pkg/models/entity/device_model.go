package entity

import (
	"github.com/google/uuid"
	"go-trans/pkg/models/consts"
)

/**

@author: victor2022
@since: 2025/1/19
*/

// DeviceInfo 设备信息
type DeviceInfo struct {
	*BaseModel
	DeviceId     string            `gorm:"index:idx_device_id" json:"device_id"`
	DeviceName   string            `gorm:"index" json:"device_name"`
	Address      string            `json:"address"`
	Port         string            `json:"port"`
	TransmitPort string            `json:"transmit_port"`
	Mode         consts.DeviceMode `gorm:"index:idx_device_mode" json:"mode"`
	PairCode     string            `json:"pair_code"`
	Connected    bool              `json:"connected"`
}

func GetNewDeviceInfo(address, port string, mode consts.DeviceMode) *DeviceInfo {
	deviceId := uuid.New().String()
	pairCode := uuid.New().String()[0:6]

	return &DeviceInfo{
		BaseModel:  GetNewBaseModel(),
		DeviceId:   deviceId,
		DeviceName: deviceId,
		Address:    address,
		Port:       port,
		PairCode:   pairCode,
		Mode:       mode,
		Connected:  false,
	}
}

func GetNewSendDeviceInfo(address, port string) *DeviceInfo {
	return GetNewDeviceInfo(address, port, consts.SendMode)
}

func GetNewReceiveDeviceInfo(deviceId, deviceName string) *DeviceInfo {
	deviceInfo := GetNewDeviceInfo("", "", consts.ReceiveMode)
	deviceInfo.DeviceId = deviceId
	deviceInfo.DeviceName = deviceName
	return deviceInfo
}

func (d *DeviceInfo) RefreshPairCode() {
	d.PairCode = uuid.New().String()[0:6]
}
