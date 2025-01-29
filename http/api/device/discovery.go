package device

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/services"
	"strings"
)

/**
  @author: victor2022
  @since: 2025/1/28
*/

func startScan(c *gin.Context) {
	ips := c.PostForm("ips")
	var ipAddrs []string
	if ips != "" {
		ipAddrs = strings.Split(ips, ",")
	} else {

	}
	services.GetServiceContext().DeviceScanService.StartScan(ipAddrs)
	response.NewSuccessResponse(c, "success")
}

func stopScan(c *gin.Context) {
	services.GetServiceContext().DeviceScanService.StopScan()
	response.NewSuccessResponse(c, "success")
}

func getScanResult(c *gin.Context) {
	results := services.GetServiceContext().DeviceScanService.GetScanResults()
	response.NewSuccessResponse(c, results)
}
