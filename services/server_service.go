package services

import (
	"go-trans/handlers"
	"log"
)

/**
  @author: victor2022
  @since: 2025/1/2
  server管理服务，用来管理接收服务器与发送服务器
*/
type ServerService struct {
	receiveServerIsOn bool
	receiveHandler    *handlers.ReceiveHandler
}

func NewServerService() *ServerService {
	return &ServerService{
		receiveServerIsOn: false,
	}
}

// StartReceiveServer startup server
func (s *ServerService) StartReceiveServer(port, basePath string) {
	receiveHandler := handlers.NewReceiveHandler(port, basePath)
	go receiveHandler.Handle()
	s.receiveServerIsOn = true
	s.receiveHandler = receiveHandler
}

// StopReceiveServer stop server
func (s *ServerService) StopReceiveServer() {
	log.Println("stopping server")
	if !s.receiveServerIsOn || nil != s.receiveHandler {
		s.receiveServerIsOn = false
	}
	s.receiveHandler.StopHandle()
	s.receiveHandler = nil
}
