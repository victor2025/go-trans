package services

/**
  @author: victor2022
  @since: 2025/1/5
*/
type ConfigService struct {
	configMap map[string]string
}

func NewConfigService() *ConfigService {

	configMap := make(map[string]string)

	configMap["http.server.port"] = "8080"
	configMap["transmit.server.port"] = "20235"
	configMap["transmit.server.file.path"] = "./received"

	return &ConfigService{
		configMap: configMap,
	}

}

func (c *ConfigService) GetOrDefault(key, defaultVal string) string {
	config, exists := c.configMap[key]
	if exists {
		return config
	}
	return defaultVal
}
