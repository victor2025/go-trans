package device

import (
	"go-trans/http/response"
	"go-trans/services"

	"github.com/gin-gonic/gin"
)

/**
  @author: victor2022
  @since: 2025/3/15
*/

func getSelfDeviceInfo(c *gin.Context) {
	selfDeviceInfo := services.GetServiceContext().DeviceService.GetSelfDeviceInfo()
	response.NewSuccessResponse(c, selfDeviceInfo)
}

// 发起配对
func pairNewDevice(c *gin.Context) {
	deviceId := c.PostForm("deviceId")
	pairCode := c.PostForm("pairCode")
	deviceScanResult := services.GetServiceContext().DeviceScanService.GetScanResultByDeviceId(deviceId)
	if deviceScanResult == nil {
		response.NewFailResponse(c, "", "device not found")
		return
	}
	err := services.GetServiceContext().DeviceService.Pair(deviceScanResult, pairCode)
	if err != nil {
		response.NewFailResponse(c, deviceId, err.Error())
		return
	}
	response.NewSuccessResponse(c, "success")
}

// 移除已配对的设备
func removePairedDevice(c *gin.Context) {
	deviceId := c.PostForm("deviceId")
	deviceInfo := services.GetServiceContext().DeviceService.GetConnectedDeviceById(deviceId)
	err := services.GetServiceContext().DeviceService.DisconnectDeviceById(deviceInfo)
	if err != nil {
		response.NewFailResponse(c, deviceId, err.Error())
		return
	}
	response.NewSuccessResponse(c, "success")
}

func refreshPairCode(c *gin.Context) {
	err := services.GetServiceContext().DeviceService.RefreshSelfPairCode()
	if err != nil {
		response.NewFailResponse(c, "refresh pair code", err.Error())
	}
	response.NewSuccessResponse(c, "success")
}
