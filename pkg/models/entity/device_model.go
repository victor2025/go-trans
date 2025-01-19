package entity

import "github.com/google/uuid"

/**
  @author: victor2022
  @since: 2025/1/19
*/
type DeviceInfo struct {
	*BaseModel
	DeviceId   string `gorm:"index:idx_device_id" json:"device_d"`
	DeviceName string `gorm:"index" json:"device_name"`
	Address    string `json:"address"`
	Port       string `json:"port"`
	PairCode   string `json:"pair_code"`
	Connected  bool   `json:"connected"`
}

func GetNewDeviceInfo(address, port string) *DeviceInfo {
	deviceId := uuid.New().String()
	pairCode := uuid.New().String()[0:6]

	return &DeviceInfo{
		BaseModel:  GetNewBaseModel(),
		DeviceId:   deviceId,
		DeviceName: deviceId,
		Address:    address,
		Port:       port,
		PairCode:   pairCode,
		Connected:  true,
	}
}
