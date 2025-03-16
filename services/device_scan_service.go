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
	"sync"
	"time"
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
	scanResult     sync.Map
	selfDeviceInfo *entity.DeviceInfo
	sn             int
	scanning       bool
	startTime      time.Time
}

func NewDeviceScanService() *DeviceScanService {
	processor := &DeviceScanService{
		sn: 0,
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

func (s *DeviceScanService) getSelfDeviceInfo() *entity.DeviceInfo {
	if s.selfDeviceInfo == nil {
		s.selfDeviceInfo = GetServiceContext().DeviceService.GetSelfDeviceInfo()
	}
	return s.selfDeviceInfo
}

// 扫描设备
func (s *DeviceScanService) scanDevice(ipAddrs []string) {
	// 初始化状态
	s.StopScan()
	s.scanning = true
	s.startTime = time.Now()
	s.scanResult = sync.Map{}
	defer func() {
		if r := recover(); r != nil {
			InfoF("scanDevice error:%s \n", r)
		}
	}()

	// 启动扫描任务
	var ips []net.IP
	if ipAddrs == nil || len(ipAddrs) == 0 {
		ips = utils.GetLocalIps()
	} else {
		for _, ipAddr := range ipAddrs {
			ips = append(ips, net.ParseIP(ipAddr))
		}
	}
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.scanDevicesByIpRange(ip)
		}()
	}
	wg.Wait()
	// 还原状态
	var count int
	s.scanResult.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	log.Printf("scanDevice finished, discover %d devices, costTime:%.2fs",
		count, time.Since(s.startTime).Seconds())
	s.scanning = false
}

func (s *DeviceScanService) scanDevicesByIpRange(localIp net.IP) {
	currSn := s.sn
	localIp = localIp.To4()
	ipCursor := make(net.IP, net.IPv4len)
	copy(ipCursor, localIp)
	ipStart := make(net.IP, net.IPv4len)
	copy(ipStart, localIp)
	ipStart[3] = 0
	ipEnd := make(net.IP, net.IPv4len)
	copy(ipEnd, ipStart)
	ipEnd[3] = 255
	// 生成ip范围
	ip3Range := make([]uint8, 2)
	if ipCursor[2] > 0 {
		ip3Range[0] = ipCursor[2] - 1
		ipStart[2] = ip3Range[0]
	}
	if ipCursor[2] < 255 {
		ip3Range[1] = ipCursor[2] + 1
		ipEnd[2] = ip3Range[1]
	}
	// 根据ip范围扫描设备
	var wg sync.WaitGroup
	for ip3 := ip3Range[0]; ip3 <= ip3Range[1]; ip3++ {
		ipCursor[2] = ip3
		for idx := 0; idx < 256; idx++ {
			ipCursor[3] = uint8(idx)
			// 判断当前任务是否要停止
			if s.sn != currSn {
				break
			}
			// 使用 WaitGroup 等待所有 scanDeviceByIP 操作完成
			wg.Add(1)
			aimedIp := ipCursor.String()
			go func() {
				defer wg.Done()
				s.scanDeviceByIp(aimedIp)
			}()
		}
	}
	wg.Wait()
	log.Printf("scanIpRange: from %s to %s finished", ipStart.String(), ipEnd.String())
}

func (s *DeviceScanService) scanDeviceByIp(ipAddr string) {
	InfoF("scanDeviceByIp for: %v\n", ipAddr)
	currSn := s.sn
	httpClient := &http.Client{}
	portStr := GetServiceContext().ConfigService.GetOrDefault(consts.HttpServerPort, consts.DefaultHttpPort)
	port, _ := strconv.Atoi(portStr)
	totalRetryCntStr := GetServiceContext().ConfigService.GetOrDefault(consts.ScanPortRetryCnt, "3")
	totalRetryCnt, _ := strconv.Atoi(totalRetryCntStr)
	for retryCnt := 0; retryCnt < totalRetryCnt; retryCnt++ {
		url := fmt.Sprintf("http://%s:%d%s", ipAddr, port, urlSuffix)
		InfoF("scanning:%v\n", url)
		request, err := http.NewRequest("GET", url, nil)
		utils.HandleError(err)
		request.Header.Set("source", consts.DeviceScanIdentityParam)
		resp, err := httpClient.Do(request)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			InfoF("receive device scan response:%v\n", string(body))
			// 解析返回值
			respEntity := &response.Entity{}
			err := json.Unmarshal(body, respEntity)
			utils.HandleError(err)
			deviceInfo := &entity.DeviceInfo{}
			response.GetStructFromResponse(respEntity, deviceInfo)
			scanInfo := dto.DeviceScanInfo{
				Ip:           ipAddr,
				Port:         strconv.Itoa(port),
				TransmitPort: deviceInfo.TransmitPort,
				DeviceId:     deviceInfo.DeviceId,
				DeviceName:   deviceInfo.DeviceName,
			}
			// 如果是本机，则不放入扫描结果中
			if scanInfo.DeviceId != s.getSelfDeviceInfo().DeviceId {
				s.saveScanResult(&scanInfo)
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

func (s *DeviceScanService) saveScanResult(scanInfo *dto.DeviceScanInfo) {
	// 查询当前设备是否已连接
	deviceInfo := GetServiceContext().DeviceService.GetConnectedDeviceById(scanInfo.DeviceId)
	if deviceInfo != nil {
		// 若设备已连接，则更新设备信息
		deviceInfo.DeviceName = scanInfo.DeviceName
		deviceInfo.Address = scanInfo.Ip
		deviceInfo.Port = scanInfo.Port
		deviceInfo.TransmitPort = scanInfo.TransmitPort
		GetServiceContext().DeviceService.UpdateDeviceById(deviceInfo)
		scanInfo.Connected = consts.YES
	}
	s.scanResult.Store(scanInfo.DeviceId, scanInfo)
}

func (s *DeviceScanService) GetScanResults() ([]*dto.DeviceScanInfo, bool) {
	results := make([]*dto.DeviceScanInfo, 0)
	s.scanResult.Range(func(key, value interface{}) bool {
		scanInfo := value.(*dto.DeviceScanInfo)
		completeConnectStatus(scanInfo)
		results = append(results, scanInfo)
		return true
	})
	return results, s.scanning
}

func (s *DeviceScanService) GetScanResultByDeviceId(deviceId string) *dto.DeviceScanInfo {
	if value, ok := s.scanResult.Load(deviceId); ok {
		scanInfo := value.(*dto.DeviceScanInfo)
		completeConnectStatus(scanInfo)
		return scanInfo
	}
	return nil
}

// 补全连接信息
func completeConnectStatus(scanInfo *dto.DeviceScanInfo) {
	deviceInfo := GetServiceContext().DeviceService.GetDeviceById(scanInfo.DeviceId)
	if deviceInfo != nil && deviceInfo.Connected {
		scanInfo.Connected = consts.YES
	} else {
		scanInfo.Connected = consts.NO
	}
}
