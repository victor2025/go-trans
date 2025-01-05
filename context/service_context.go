package context

import (
	"go-trans/services"
	"sync"
)

/**
  @author: victor2022
  @since: 2025/1/5
*/

type SystemContext struct {
	serverService *services.TransmitService
	configService *services.ConfigService
}

var systemContext *SystemContext
var once = sync.Once{}

// GetSystemContext 获取系统上下文
func GetSystemContext() *SystemContext {
	once.Do(func() {
		systemContext = &SystemContext{
			serverService: services.NewServerService(),
			configService: services.NewConfigService(),
		}
	})
	return systemContext
}

func (s *SystemContext) StartReceiveServer() {
	port := s.configService.GetOrDefault("transmit.server.port", "20235")
	filePath := s.configService.GetOrDefault("transmit.server.filepath", "./received")
	s.serverService.StartReceiveServer(port, filePath)
}

func (s *SystemContext) StopReceiveServer() {
	s.serverService.StopReceiveServer()
}
