package services

import (
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
	"strings"
)

/**
@author: victor2022
@since: 2025/1/5
*/

const (
	baseConfigPath = "./res/config/config.json"
)

type ConfigService struct {
	configMap map[string]any
	db        *gorm.DB
}

func NewConfigService() *ConfigService {
	configMap := make(map[string]any)
	configService := &ConfigService{
		configMap: configMap,
	}
	configService.loadConfigFromFiles()
	return configService
}

func (c *ConfigService) loadConfigFromFiles() {
	err := utils.LoadJsonFile(baseConfigPath, &c.configMap)
	utils.HandleError(err)
}

// GetOrDefaultFromFile 从文件中获取配置
func (c *ConfigService) GetOrDefaultFromFile(key, defaultVal string) string {
	parts := strings.Split(key, ".")
	config := defaultVal
	currResult := c.configMap
	for idx, part := range parts {
		if value, ok := currResult[part].(map[string]interface{}); ok {
			currResult = value
		} else if idx == len(parts)-1 {
			finalVal := currResult[part]
			if finalVal != nil {
				config = finalVal.(string)
			}
		}
	}
	return config
}

func (c *ConfigService) GetOrDefault(key, defaultVal string) string {
	var config string
	// 先从db取
	var configInfo *entity.ConfigInfo
	tx := c.db.Find(&configInfo, "key = ?", key)
	if tx.RowsAffected == 0 {
		// db中不存在，则从本地文件取
		config = c.GetOrDefaultFromFile(key, defaultVal)
		// 保存到db
		configInfo = entity.GetNewConfigInfo(key, config)
		c.db.Save(&configInfo)
	}
	return configInfo.ConfigValue
}
