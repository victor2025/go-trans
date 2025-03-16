package device

import (
	"encoding/json"
	"go-trans/http/response"
	"go-trans/pkg/models/consts"
	"go-trans/pkg/models/entity"
	"go-trans/services"

	"github.com/gin-gonic/gin"
)

/**
@author: victor2022
@since: 2025/1/27
	内部接口
*/

// 用于响应扫描请求，返回本机设备id
func ping(c *gin.Context) {
	validateParam(c)
	selfDeviceInfo := services.GetServiceContext().DeviceService.GetSelfDeviceInfo()
	returnObj := map[string]string{
		"deviceId":     selfDeviceInfo.DeviceId,
		"deviceName":   selfDeviceInfo.DeviceName,
		"transmitPort": selfDeviceInfo.TransmitPort,
	}
	response.NewSuccessResponse(c, returnObj)
}

// 用于配对
func pair(c *gin.Context) {
	validateParam(c)
	// 获取请求者的ip地址
	clientIp := c.ClientIP()
	if clientIp == "" {
		response.NewFailResponse(c, "", "无法获取客户端IP地址")
		return
	}
	deviceInfoStr := c.PostForm("deviceInfo")
	var deviceInfo *entity.DeviceInfo
	err := json.Unmarshal([]byte(deviceInfoStr), &deviceInfo)
	pairCode := c.PostForm("pairCode")
	if err != nil {
		response.NewFailResponse(c, "", "invalid device info: "+err.Error())
		return
	}
	// 设置clientIP
	deviceInfo.Address = clientIp
	err = services.GetServiceContext().DeviceService.BePaired(deviceInfo, pairCode)
	if err != nil {
		response.NewFailResponse(c, "", err.Error())
		return
	}
	// 返回port
	transmitPort := services.GetServiceContext().ConfigService.GetOrDefault(consts.TransServerPort, consts.DefaultHttpPort)
	selfDeviceInfo := services.GetServiceContext().DeviceService.GetSelfDeviceInfo()
	response.NewSuccessResponse(c, &map[string]string{
		"deviceId":     selfDeviceInfo.DeviceId,
		"transmitPort": transmitPort,
	})
}

func validateParam(c *gin.Context) {
	source := c.GetHeader("source")
	if source != consts.DeviceScanIdentityParam {
		response.NewFailResponse(c, nil, "invalid source: "+source)
		panic("invalid source: " + source)
	}
}
