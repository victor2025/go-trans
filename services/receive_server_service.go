package services

import (
	"fmt"
	"go-trans/pkg/transmit/handlers"
	"go-trans/utils"
	"log"
	"strconv"
)

/*
*

	@author: victor2022
	@since: 2025/1/2
	传输管理服务，用来管理接收服务器与发送服务器
*/
type TransmitService struct {
	receiveServerIsOn bool
	receiveHandler    *handlers.ReceiveHandler
	ServerPort        int
}

func NewServerService() *TransmitService {
	return &TransmitService{
		receiveServerIsOn: false,
	}
}

// StartReceiveServer startup server
func (s *TransmitService) StartReceiveServer(port, basePath string) {
	var err error
	s.ServerPort, err = strconv.Atoi(port)
	utils.HandleError(err, utils.PanicOnError)
	if !s.receiveServerIsOn {
		var receiveHandler *handlers.ReceiveHandler
		for retryCnt := 0; retryCnt < 10; retryCnt++ {
			receiveHandler = handlers.NewReceiveHandler(fmt.Sprintf("%d", s.ServerPort), basePath)
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("start transmit server panic: %v", r)
						err = fmt.Errorf("%v", r)
					}
				}()
				receiveHandler.Handle()
			}()
			if err == nil {
				break
			}
			s.ServerPort += 1
		}
		utils.HandleError(err, utils.PanicOnError)
		s.receiveServerIsOn = true
		s.receiveHandler = receiveHandler
	} else {
		log.Printf("transmit server is already on, port:%d", s.ServerPort)
	}
}

// StopReceiveServer stop server
func (s *TransmitService) StopReceiveServer() {
	log.Println("stopping transmit server")
	if s.receiveServerIsOn {
		s.receiveServerIsOn = false
	}
	if s.receiveHandler != nil {
		s.receiveHandler.StopHandle()
	}
	s.receiveHandler = nil
	log.Println("transmit server stopped")
}
