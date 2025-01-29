package main

import (
	"go-trans/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

/**
  @author: victor2022
  @since: 2025/1/2
*/

func main() {
	// 启动应用
	app.Run()

	// 监听 SIGINT (Ctrl+C) 和 SIGTERM (kill 命令)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Println("receive shutdown signal:", sig)
	// 关闭应用
	app.Shutdown()

}
