package device

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/pkg/models/consts"
	"go-trans/services"
)

/*
*

	@author: victor2022
	@since: 2025/1/27
*/
func ping(c *gin.Context) {
	source := c.Query("source")
	if source != consts.DeviceScanIdentityParam {
		response.NewFailResponse(c, nil, "invalid source")
		return
	}
	selfDeviceInfo := services.GetServiceContext().DeviceService.GetSelfDeviceInfo()
	response.NewSuccessResponse(c, selfDeviceInfo.DeviceId)
}
