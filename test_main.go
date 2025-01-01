package main

import (
	"go-trans/services"
	"log"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/2
*/

func main() {
	serverService := services.NewServerService()
	serverService.StartReceiveServer("20235", ".received")
	for i := 0; i < 2; i++ {
		time.Sleep(1 * time.Second)
		log.Printf("wait: %d", i)
	}
	serverService.StopReceiveServer()
	for {
	}
}
