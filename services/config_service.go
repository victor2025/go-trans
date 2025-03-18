package services

import (
	"fmt"
	"go-trans/pkg/models/entity"
	"go-trans/utils"
	"gorm.io/gorm"
	"strings"
	"sync"
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
	startupConfig map[string]any
	configCache   sync.Map
	db            *gorm.DB
}

func NewConfigService(config map[string]any) *ConfigService {
	startupConfig := make(map[string]any)
	configService := &ConfigService{
		startupConfig: startupConfig,
		configCache:   sync.Map{},
	}
	// 从入参中读取配置
	if config != nil {
		for k, v := range config {
			startupConfig[k] = v
			configService.saveStartupConfigToCache(k, v)
		}
	}
	return configService
}

func (c *ConfigService) loadConfigFromFiles() {
	err := utils.LoadJsonFile(baseConfigPath, &c.startupConfig)
	utils.HandleError(err)
}

// GetOrDefaultFromStartupConfig 获取启动配置
func (c *ConfigService) GetOrDefaultFromStartupConfig(key, defaultVal string) string {
	parts := strings.Split(key, ".")
	config := defaultVal
	currResult := c.startupConfig
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

// 将启动参数保存到缓存中
func (c *ConfigService) saveStartupConfigToCache(currKey string, config any) {
	if valMap, ok := config.(map[string]interface{}); ok {
		for key, val := range valMap {
			c.saveStartupConfigToCache(currKey+"."+key, val)
		}
	} else {
		c.configCache.Store(currKey, config)
	}
}

func (c *ConfigService) GetOrDefault(key, defaultVal string) string {
	if val, _ := c.configCache.Load(key); val != nil {
		return fmt.Sprint(val)
	}
	if c.db == nil {
		return defaultVal
	}
	// 先从db取
	var configInfo *entity.ConfigInfo
	tx := c.db.Find(&configInfo, "config_id = ?", key)
	if tx.RowsAffected == 0 {
		// 默认配置
		configInfo = entity.GetNewConfigInfo(key, defaultVal)
	}
	c.configCache.Store(key, configInfo.ConfigValue)
	return configInfo.ConfigValue
}

// SetTempConfig 设置临时参数
func (c *ConfigService) SetTempConfig(key, val string) {
	c.configCache.Store(key, val)
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
	c.configCache.Store(key, configInfo.ConfigValue)
	return nil
}
