package device

import "github.com/gin-gonic/gin"

/*
*

	@author: victor2022
	@since: 2025/1/27
*/
func RegisterRouter(router *gin.RouterGroup) {
	// 内部接口
	router.GET("inner/ping", ping)
	router.POST("inner/pair", pair)

	// 设备发现
	router.POST("discovery/startScan", startScan)
	router.POST("discovery/stopScan", stopScan)
	router.GET("discovery/getScanResult", getScanResult)
}
