package system

import (
	"go-trans/http/response"
	"go-trans/services"

	"github.com/gin-gonic/gin"
)

/**
  @author: victor2022
  @since: 2025/1/28
*/

// 更新设置
func updateSetting(c *gin.Context) {
	key := c.PostForm("key")
	config := c.PostForm("config")
	if key == "" || config == "" {
		panic("key or value is empty")
	}
	err := services.GetServiceContext().ConfigService.UpdateConfig(key, config)
	if err != nil {
		response.NewFailResponse(c, "", err.Error())
		return
	}
	response.NewSuccessResponse(c, map[string]string{
		"key":    key,
		"config": config,
	})
}

// 更新设置
func getSetting(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		panic("key or value is empty")
	}
	config := services.GetServiceContext().ConfigService.GetOrDefault(key, "")
	response.NewSuccessResponse(c, map[string]string{
		"key":    key,
		"config": config,
	})
}
