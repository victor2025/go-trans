package device

import (
	"github.com/gin-gonic/gin"
	"go-trans/http/response"
	"go-trans/services"
)

/**
  @author: victor2022
  @since: 2025/1/28
*/

func startScan(c *gin.Context) {
	services.GetServiceContext().DeviceScanService.StartScan()
	response.NewSuccessResponse(c, "")
}

func getScanResult(c *gin.Context) {
	results := services.GetServiceContext().DeviceScanService.GetScanResults()
	response.NewSuccessResponse(c, results)
}
