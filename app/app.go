package app

import (
	"go-trans/http"
	"log"
	"time"
)

/**
  @author: victor2022
  @since: 2025/1/24
*/
func Startup() {
	// 启动httpServer
	start := time.Now()
	go http.StartHttpServer()
	log.Printf("go-trans-core startup in %dms\n", time.Since(start).Microseconds())
}

func Shutdown() {
	// todo shutdown
}
