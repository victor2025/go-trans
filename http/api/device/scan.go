package device

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/pkg/models/consts"
	"go-trans/services"
)

/**

@author: victor2022
@since: 2025/1/27
*/

// 用于响应扫描请求，返回本机设备id
func ping(c *gin.Context) {
	validateParam(c)
	selfDeviceInfo := services.GetServiceContext().DeviceService.GetSelfDeviceInfo()
	returnObj := map[string]string{
		"deviceId":   selfDeviceInfo.DeviceId,
		"deviceName": selfDeviceInfo.DeviceName,
	}
	response.NewSuccessResponse(c, returnObj)
}

// 用于配对
func pair(c *gin.Context) {
	validateParam(c)
	deviceId := c.PostForm("deviceId")
	deviceName := c.PostForm("deviceName")
	pairCode := c.PostForm("pairCode")
	err := services.GetServiceContext().DeviceService.PairDeviceForReceive(deviceId, deviceName, pairCode)
	if err != nil {
		response.NewFailResponse(c, "", err.Error())
		return
	}
	response.NewSuccessResponse(c, "success")
}

func validateParam(c *gin.Context) {
	source := c.GetHeader("source")
	if source != consts.DeviceScanIdentityParam {
		response.NewFailResponse(c, nil, "invalid source: "+source)
		panic("invalid source: " + source)
	}
}
