package services

import (
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
	DB                   *gorm.DB
	BusService           *BusService
}

var serviceContext *ServiceContext
var once = sync.Once{}

// GetServiceContext 获取系统上下文
func GetServiceContext() *ServiceContext {
	once.Do(func() {
		// 配置服务
		configService := NewConfigService()
		// 持久层配置
		dbLocation := configService.GetOrDefault("orm.db.location", "./res/db/dev.db")
		DB := orm.GetDb(dbLocation)
		// 总线配置
		busTopic := configService.GetOrDefault("msg.bus.topic", "msg.bus")

		serviceContext = &ServiceContext{
			ReceiveServerService: NewServerService(),
			ConfigService:        configService,
			SendTaskService:      NewTaskService(DB),
			BusService:           NewBusService(busTopic),
			DeviceService:        NewDeviceService(DB),
			DB:                   DB,
		}
	})
	return serviceContext
}

func (s *ServiceContext) StartReceiveServer() {
	port := s.ConfigService.GetOrDefault("transmit.server.port", "20235")
	filePath := s.ConfigService.GetOrDefault("transmit.server.filepath", "./received")
	s.ReceiveServerService.StartReceiveServer(port, filePath)
}

func (s *ServiceContext) StopReceiveServer() {
	s.ReceiveServerService.StopReceiveServer()
}
