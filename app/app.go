package app

import (
	"go-trans/http"
)

/**
  @author: victor2022
  @since: 2025/1/24
*/
func Run() {
	// 初始化系统
	RegisterMsgConsumers()
	RegisterProcessors()
	// 启动http服务
	http.StartHttpServer()
}

func Shutdown() {
	// todo shutdown
}
