package config

import (
	logger2 "github.com/grobza/proxyless-grpc-lb/hello-world/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ReadConfig reads the config data from file
func ReadConfig() (*viper.Viper, error) {
	logger2.Logger.Debug("Reading configuration", zap.String("file", "/var/run/config/app.yaml"))
	v := viper.New()
	v.SetConfigFile("/var/run/config/app.yaml")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err

}
