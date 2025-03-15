package app

import (
	"go-trans/http"
	"go-trans/services"
)

/**
  @author: victor2022
  @since: 2025/1/27
*/

// StartupServer 启动服务器
func StartupServer() int {
	// 启动接收服务器
	services.GetServiceContext().StartReceiveServer()
	// 启动http服务器
	return http.StartHttpServer()
}

// ShutdownServer 关闭服务器
func ShutdownServer() {
	// 关闭http服务器
	http.ShutdownHttpServer()
	// 关闭接收服务器
	services.GetServiceContext().StopReceiveServer()
}
