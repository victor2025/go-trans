package services

import (
	"go-trans/pkg/transmit/handlers"
	"log"
)

/**
  @author: victor2022
  @since: 2025/1/2
  传输管理服务，用来管理接收服务器与发送服务器
*/
type TransmitService struct {
	receiveServerIsOn bool
	receiveHandler    *handlers.ReceiveHandler
}

func NewServerService() *TransmitService {
	return &TransmitService{
		receiveServerIsOn: false,
	}
}

// StartReceiveServer startup server
func (s *TransmitService) StartReceiveServer(port, basePath string) {
	receiveHandler := handlers.NewReceiveHandler(port, basePath)
	go receiveHandler.Handle()
	s.receiveServerIsOn = true
	s.receiveHandler = receiveHandler
}

// StopReceiveServer stop server
func (s *TransmitService) StopReceiveServer() {
	log.Println("Stopping server")
	if !s.receiveServerIsOn || nil != s.receiveHandler {
		s.receiveServerIsOn = false
	}
	s.receiveHandler.StopHandle()
	s.receiveHandler = nil
}
