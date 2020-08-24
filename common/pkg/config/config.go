package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	logger "github.com/asishrs/proxyless-grpc-lb/common/pkg/logger"
)

// ReadConfig reads the config data from file
func ReadConfig() (*viper.Viper, error) {
	logger.Logger.Debug("Reading configuration", zap.String("file", "/var/run/config/app.yaml"))
	v := viper.New()
	v.SetConfigFile("/var/run/config/app.yaml")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err

}
