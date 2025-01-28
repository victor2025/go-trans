package entity

/**
@author: victor2022
@since: 2025/1/28
*/

// ConfigInfo 配置信息
type ConfigInfo struct {
	*BaseModel
	ConfigId    string `json:"config_id"`
	ConfigValue string `json:"config_value"`
}

func GetNewConfigInfo(configId, configVal string) *ConfigInfo {
	return &ConfigInfo{
		BaseModel:   GetNewBaseModel(),
		ConfigId:    configId,
		ConfigValue: configVal,
	}
}
