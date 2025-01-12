package context

import (
	"go-trans/services"
	"sync"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

type ServiceContext struct {
	serverService *services.TransmitService
	configService *services.ConfigService
}

var serviceContext *ServiceContext
var once = sync.Once{}

// GetServiceContext 获取系统上下文
func GetServiceContext() *ServiceContext {
	once.Do(func() {
		serviceContext = &ServiceContext{
			serverService: services.NewServerService(),
			configService: services.NewConfigService(),
		}
	})
	return serviceContext
}

func (s *ServiceContext) StartReceiveServer() {
	port := s.configService.GetOrDefault("transmit.server.port", "20235")
	filePath := s.configService.GetOrDefault("transmit.server.filepath", "./received")
	s.serverService.StartReceiveServer(port, filePath)
}

func (s *ServiceContext) StopReceiveServer() {
	s.serverService.StopReceiveServer()
}

func (s *ServiceContext) GetConfigOrDefault(key, defaultVal string) string {
	return s.configService.GetOrDefault(key, defaultVal)
}
