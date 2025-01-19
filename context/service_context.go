package context

import (
	"go-trans/pkg/orm"
	"go-trans/services"
	"gorm.io/gorm"
	"sync"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

type ServiceContext struct {
	ServerService *services.TransmitService
	ConfigService *services.ConfigService
	TaskService   *services.TaskService
	DB            *gorm.DB
	BusService    *services.BusService
}

var serviceContext *ServiceContext
var once = sync.Once{}

// GetServiceContext 获取系统上下文
func GetServiceContext() *ServiceContext {
	once.Do(func() {
		// 配置服务
		configService := services.NewConfigService()
		// 持久层配置
		dbLocation := configService.GetOrDefault("orm.db.location", "./res/db/dev.db")
		DB := orm.GetDb(dbLocation)
		// 总线配置
		busTopic := configService.GetOrDefault("msg.bus.topic", "msg.bus")

		serviceContext = &ServiceContext{
			ServerService: services.NewServerService(),
			ConfigService: configService,
			TaskService:   services.NewTaskService(DB),
			BusService:    services.NewBusService(busTopic),
			DB:            DB,
		}
	})
	return serviceContext
}

func (s *ServiceContext) StartReceiveServer() {
	port := s.ConfigService.GetOrDefault("transmit.server.port", "20235")
	filePath := s.ConfigService.GetOrDefault("transmit.server.filepath", "./received")
	s.ServerService.StartReceiveServer(port, filePath)
}

func (s *ServiceContext) StopReceiveServer() {
	s.ServerService.StopReceiveServer()
}
