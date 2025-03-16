package dto

/**
  @author: victor2022
  @since: 2025/1/27
*/

// DeviceScanInfo 设备扫描信息
type DeviceScanInfo struct {
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	TransmitPort string `json:"transmitPort"`
	DeviceId     string `json:"deviceId"`
	DeviceName   string `json:"deviceName"`
	Connected    string `json:"connected"`
}
