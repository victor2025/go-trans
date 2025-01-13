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
}

var serviceContext *ServiceContext
var once = sync.Once{}

// GetServiceContext 获取系统上下文
func GetServiceContext() *ServiceContext {
	once.Do(func() {
		configService := services.NewConfigService()
		dbLocation := configService.GetOrDefault("orm.db.location", "./res/db/dev.db")
		DB := orm.GetDb(dbLocation)
		serviceContext = &ServiceContext{
			ServerService: services.NewServerService(),
			ConfigService: configService,
			TaskService:   services.NewTaskService(DB),
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
