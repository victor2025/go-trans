package app

/*
*

	@author: victor2022
	@since: 2025/1/24
*/
func Run() {
	// 初始化系统
	RegisterMsgConsumers()
	RegisterProcessors()
	// 启动服务器
	StartupServer()
}

func Shutdown() {
	// todo shutdown
}
