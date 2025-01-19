package services

import (
	"go-trans/utils"
	"log"
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
	if err != nil {
		log.Println(err)
	}

}

func (c *ConfigService) GetOrDefault(key, defaultVal string) string {
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
