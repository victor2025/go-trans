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
	"log"
	"net"
	"net/http"
	"strconv"
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
	scanResult     map[string]dto.DeviceScanInfo
	selfDeviceInfo *entity.DeviceInfo
	sn         int
}

func NewDeviceScanService(info *entity.DeviceInfo) *DeviceScanService {
	processor := &DeviceScanService{
		scanResult:     make(map[string]dto.DeviceScanInfo),
		selfDeviceInfo: info,
		sn:         0,
	}
	return processor
}

// StartScan 启动新线程扫描
func (s *DeviceScanService) StartScan(ipAddrs []string) {
	go s.scanDevice(ipAddrs)
}

// StopScan 关闭正在进行的扫描任务
func (s *DeviceScanService) StopScan() {
	s.sn++
}

// 扫描设备
func (s *DeviceScanService) scanDevice(ipAddrs []string) {
	s.StopScan()
	s.scanResult = make(map[string]dto.DeviceScanInfo)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("scanDevice error:%s \n", r)
		}
	}()

	var ips []net.IP
	if ipAddrs == nil || len(ipAddrs) == 0 {
		ips = utils.GetLocalIps()
	} else {
		for _, ipAddr := range ipAddrs {
			ips = append(ips, net.ParseIP(ipAddr))
		}
	}
	for _, ip := range ips {
		go s.scanDevicesByIpRange(ip)
	}

}

func (s *DeviceScanService) scanDevicesByIpRange(localIp net.IP) {
	currSn := s.sn
	localIp = localIp.To4()
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
			// 判断当前任务是否要停止
			if s.sn != currSn {
				break
			}
			// 扫描
			go s.scanDeviceByIp(ipCursor.String())
		}
	}
}

func (s *DeviceScanService) scanDeviceByIp(ipAddr string) {
	currSn := s.sn
	httpClient := &http.Client{}
	portStr := GetServiceContext().ConfigService.GetOrDefault(consts.HttpServerPort, "9210")
	port, _ := strconv.Atoi(portStr)
	for retryCnt := 0; retryCnt < 10; retryCnt++ {
		url := fmt.Sprintf("http://%s:%d%s", ipAddr, port, urlSuffix)
		request, err := http.NewRequest("GET", url, nil)
		utils.HandleError(err)
		request.Header.Set("source", consts.DeviceScanIdentityParam)
		resp, err := httpClient.Do(request)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("receive device scan response:%v\n", string(body))
			// 解析返回值
			respEntity := &response.Entity{}
			err := json.Unmarshal(body, respEntity)
			utils.HandleError(err)
			contentMap := respEntity.Content.(map[string]interface{})
			scanInfo := dto.DeviceScanInfo{
				Ip:       ipAddr,
				Port:     strconv.Itoa(port),
				DeviceId: contentMap["deviceId"].(string),
			}
			// 如果是本机，则不放入扫描结果中
			if scanInfo.DeviceId != s.selfDeviceInfo.DeviceId {
				s.scanResult[scanInfo.DeviceId] = scanInfo
				log.Printf("discover new device, info:%s \n", scanInfo)
			}
		}
		port++

		// 判断是否发起了新请求
		if s.sn != currSn {
			break
		}
	}

}

func (s *DeviceScanService) GetScanResults() []*dto.DeviceScanInfo {
	results := make([]*dto.DeviceScanInfo, 0)
	for _, scanInfo := range s.scanResult {
		results = append(results, &scanInfo)
	}
	return results
}
