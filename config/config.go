package config

import "encoding/json"

import (
	"github.com/winc-link/hummingbird-sdk-go/service"
)

var baseConfig *BaseConfig

type BaseConfig struct {
}

func InitConfig(sd *service.DriverService) {
	customParam := sd.GetCustomParam()
	baseConfig = &BaseConfig{}
	if customParam != "" {
		err := json.Unmarshal([]byte(customParam), &baseConfig)
		if err != nil {
			sd.GetLogger().Error(err)
		}
	}

}

func GetConfig() *BaseConfig {
	return baseConfig
}
