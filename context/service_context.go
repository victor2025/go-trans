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
}

var systemContext *SystemContext
var once = sync.Once{}

// GetSystemContext 获取系统上下文
func GetSystemContext() *SystemContext {
	once.Do(func() {
		systemContext = &SystemContext{
			serverService: services.NewServerService(),
		}
	})
	return systemContext
}

func (s *SystemContext) StartReceiveServer() {
	// todo load props from propService
	s.serverService.StartReceiveServer("20235", "./.received")
}

func (s *SystemContext) StopReceiveServer() {
	s.serverService.StopReceiveServer()
}
