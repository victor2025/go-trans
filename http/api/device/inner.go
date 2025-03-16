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
		"deviceId":   selfDeviceInfo.DeviceId,
		"deviceName": selfDeviceInfo.DeviceName,
	}
	response.NewSuccessResponse(c, returnObj)
}

// 用于配对
func pair(c *gin.Context) {
	validateParam(c)
	deviceInfoStr := c.PostForm("deviceInfo")
	var deviceInfo *entity.DeviceInfo
	err := json.Unmarshal([]byte(deviceInfoStr), &deviceInfo)
	pairCode := c.PostForm("pairCode")
	if err != nil {
		response.NewFailResponse(c, "", "invalid device info: "+err.Error())
		return
	}
	err = services.GetServiceContext().DeviceService.BePaired(deviceInfo, pairCode)
	if err != nil {
		response.NewFailResponse(c, "", err.Error())
		return
	}
	// 返回port
	transServerPort := services.GetServiceContext().ConfigService.GetOrDefault(consts.TransServerPort, "20235")
	selfDeviceInfo := services.GetServiceContext().DeviceService.GetSelfDeviceInfo()
	response.NewSuccessResponse(c, &map[string]string{
		"deviceId":        selfDeviceInfo.DeviceId,
		"transServerPort": transServerPort,
	})
}

func validateParam(c *gin.Context) {
	source := c.GetHeader("source")
	if source != consts.DeviceScanIdentityParam {
		response.NewFailResponse(c, nil, "invalid source: "+source)
		panic("invalid source: " + source)
	}
}
