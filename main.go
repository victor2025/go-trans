package main

import (
	"encoding/json"
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
	configStr := "{\"orm\":{\"db\":{\"location\":\"/home/victor2022/.local/share/go_trans_flutter/db/app.db\"}},\"transmit\":{\"server\":{\"filepath\":\"/home/victor2022/下载/gotrans\"}}}"
	var config map[string]any
	json.Unmarshal([]byte(configStr), &config)
	app.RunWithConfig(config)

	// 监听 SIGINT (Ctrl+C) 和 SIGTERM (kill 命令)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Println("receive shutdown signal:", sig)
	// 关闭应用
	app.Shutdown()

}
