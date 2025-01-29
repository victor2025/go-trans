package app

import "go-trans/services"

/**
@author: victor2022
@since: 2025/1/24
*/

// RunWithConfig 带配置启动主应用
func RunWithConfig(config map[string]any) {
	services.InitServiceContext(config)
	Run()
}

// Run 启动主应用
func Run() {
	// 初始化系统
	RegisterMsgConsumers()
	RegisterProcessors()
	// 启动服务器
	StartupServer()
}

// Shutdown 关闭主应用
func Shutdown() {
	// 关闭服务器
	ShutdownServer()
	// 关闭所有处理器
	UnregisterProcessors()
}
