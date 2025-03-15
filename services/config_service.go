package services

import (
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

/**
@author: victor2022
@since: 2025/1/5
*/

const (
	baseConfigPath = "./res/config/config.json"
)

type ConfigService struct {
	configMap   map[string]any
	configCache map[string]string
	db          *gorm.DB
}

func NewConfigService(config map[string]any) *ConfigService {
	configMap := make(map[string]any)
	configService := &ConfigService{
		configMap:   configMap,
		configCache: make(map[string]string),
	}
	// 从入参中读取配置
	if config != nil {
		for k, v := range config {
			configMap[k] = v
		}
	}
	return configService
}

func (c *ConfigService) loadConfigFromFiles() {
	err := utils.LoadJsonFile(baseConfigPath, &c.configMap)
	utils.HandleError(err)
}

// GetOrDefaultFromMap 从map中获取配置，已废弃
func (c *ConfigService) GetOrDefaultFromMap(key, defaultVal string) string {
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
	if c.configCache[key] != "" {
		return c.configCache[key]
	}
	// 先从db取
	var configInfo *entity.ConfigInfo
	tx := c.db.Find(&configInfo, "config_id = ?", key)
	if tx.RowsAffected == 0 {
		// 默认配置
		configInfo = entity.GetNewConfigInfo(key, defaultVal)
	}
	c.configCache[key] = defaultVal
	return configInfo.ConfigValue
}

// UpdateConfig 更新配置
func (c *ConfigService) UpdateConfig(key, config string) error {
	var configInfo *entity.ConfigInfo
	tx := c.db.Find(&configInfo, "config_id = ?", key)
	if tx.RowsAffected == 0 {
		configInfo = entity.GetNewConfigInfo(key, config)
	} else {
		configInfo.ConfigValue = config
		configInfo.GmtModify = time.Now()
	}
	err := c.db.Save(&configInfo).Error
	if err != nil {
		return err
	}
	c.configCache[key] = configInfo.ConfigValue
	return nil
}
