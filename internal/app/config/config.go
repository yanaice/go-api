package config

import "time"

var Conf Config

type Config struct {
	Server struct {
		Port                 uint          `yaml:"port"`
		ReadHeaderTimeout    time.Duration `yaml:"read-header-timeout"`
		WriteTimeout         time.Duration `yaml:"write-timeout"`
		GracefulShutdownTime time.Duration `yaml:"graceful-shutdown-time"`
	} `yaml:"server"`
	MongoDB struct {
		URI      string `yaml:"uri"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Schema   string `yaml:"schema"`
	} `yaml:"mongodb"`
	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	BanTokenGuardThreshold  int64 `yaml:"ban-token-guard-threshold"`
	AllowConcurrentSessions bool  `yaml:"allow-concurrent-sessions"`
}

const ApplicationName = "petstore"
