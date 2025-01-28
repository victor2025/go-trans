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

const (
	urlSuffix = "/device/inner/ping"
)

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
				s.scanDevicesByIpRange(ipnet.IP.To4())
			}
		}
	}

}

func (s *DeviceScanService) scanDevicesByIpRange(localIp net.IP) {
	localIp = localIp.To4()
	localIpAddr := localIp.String()
	ipCursor := localIp
	// 生成ip范围
	ip3Range := make([]uint8, 2)
	if ipCursor[2] > 0 {
		ip3Range[0] = ipCursor[2] - 1
	}
	if ipCursor[2] < 255 {
		ip3Range[1] = ipCursor[2] + 1
	}
	// 根据ip范围扫描设备
	for ip3 := ip3Range[0]; ip3 <= ip3Range[1]; ip3++ {
		ipCursor[2] = ip3
		for idx := 0; idx < 256; idx++ {
			ipCursor[3] = uint8(idx)
			ipAddr := ipCursor.To4().String()
			if ipAddr == localIpAddr {
				continue
			}
			// 扫描
		}
	}
}

func (s *DeviceScanService) scanDeviceByIp(ipAddr string) {

}

func (s *DeviceScanService) GetScanResults() []*dto.DeviceScanInfo {
	results := make([]*dto.DeviceScanInfo, len(s.scanResult))
	for _, scanInfo := range s.scanResult {
		results = append(results, &scanInfo)
	}
	return results
}
