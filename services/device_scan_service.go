package services

import (
	"go-trans/pkg/models/dto"
	"go-trans/utils"
	"log"
	"net"
)

/**
  @author: victor2022
  @since: 2025/1/28
	设备发现服务
*/

// DeviceScanService 设备扫描处理器
type DeviceScanService struct {
	scanResult map[string]dto.DeviceScanInfo
}

func NewDeviceScanService() *DeviceScanService {
	processor := &DeviceScanService{
		scanResult: make(map[string]dto.DeviceScanInfo),
	}
	return processor
}

// StartScan 启动新线程扫描
func (s *DeviceScanService) StartScan() {
	go s.scanDevice()
}

// 扫描设备
func (s *DeviceScanService) scanDevice() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("scanDevice error:%s \n", r)
		}
	}()
	addrs, err := net.InterfaceAddrs()
	utils.HandleError(err, utils.PanicOnError)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addrStr := ipnet.IP.String()
				if len(addrStr) > 0 {
					// todo 扫描
				}
			}
		}
	}

}

func (s *DeviceScanService) scanDevicesByIpRange(myIp string) {
	// todo 根据ip范围扫描设备
}

func (s *DeviceScanService) GetScanResults() []*dto.DeviceScanInfo {
	results := make([]*dto.DeviceScanInfo, len(s.scanResult))
	for _, scanInfo := range s.scanResult {
		results = append(results, &scanInfo)
	}
	return results
}
