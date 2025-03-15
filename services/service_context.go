package services

import (
	"go-trans/pkg/models/consts"
	"go-trans/pkg/orm"
	"gorm.io/gorm"
	"sync"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

type ServiceContext struct {
	ReceiveServerService *TransmitService
	ConfigService        *ConfigService
	SendTaskService      *SendTaskService
	DeviceService        *DeviceService
	DeviceScanService    *DeviceScanService
	DB                   *gorm.DB
	BusService           *BusService
}

var serviceContext *ServiceContext
var once = sync.Once{}

func InitServiceContext(config map[string]any) {
	once.Do(func() {
		// 配置服务
		configService := NewConfigService(config)
		// 持久层配置
		dbLocation := configService.GetOrDefault(consts.OrmDbLocation, "./res/db/app.db")
		DB := orm.GetDb(dbLocation)
		// 为configService配置db
		configService.db = DB
		// 总线配置
		busTopic := configService.GetOrDefault(consts.MsgBusTopic, "msg.bus")

		deviceService := NewDeviceService(DB)

		serviceContext = &ServiceContext{
			ReceiveServerService: NewServerService(),
			ConfigService:        configService,
			SendTaskService:      NewTaskService(DB),
			BusService:           NewBusService(busTopic),
			DeviceService:        deviceService,
			DeviceScanService:    NewDeviceScanService(deviceService.GetSelfDeviceInfo()),
			DB:                   DB,
		}
	})
}

// GetServiceContext 获取系统上下文
func GetServiceContext() *ServiceContext {
	if serviceContext == nil {
		InitServiceContext(nil)
	}
	return serviceContext
}

func (s *ServiceContext) StartReceiveServer() {
	port := s.ConfigService.GetOrDefault(consts.TransServerPort, "20235")
	filePath := s.ConfigService.GetOrDefault(consts.TransServerFilepath, "./received")
	go s.ReceiveServerService.StartReceiveServer(port, filePath)
}

func (s *ServiceContext) StopReceiveServer() {
	s.ReceiveServerService.StopReceiveServer()
}
