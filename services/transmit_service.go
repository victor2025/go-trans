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
	if !s.receiveServerIsOn {
		receiveHandler := handlers.NewReceiveHandler(port, basePath)
		go receiveHandler.Handle()
		s.receiveServerIsOn = true
		s.receiveHandler = receiveHandler
	} else {
		log.Printf("Transmit server is already on, port:%s", port)
	}
}

// StopReceiveServer stop server
func (s *TransmitService) StopReceiveServer() {
	log.Println("Stopping server")
	if !s.receiveServerIsOn {
		s.receiveServerIsOn = false
	}
	if s.receiveHandler != nil {
		s.receiveHandler.StopHandle()
	}
	s.receiveHandler = nil
}
