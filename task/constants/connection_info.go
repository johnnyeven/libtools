package constants

import "github.com/johnnyeven/libtools/conf/presets"

type ConnectionInfo struct {
	Protocol string           `conf:"env"`
	Host     string           `conf:"env,upstream"`
	Port     int32            `conf:"env"`
	UserName string           `conf:"env"`
	Password presets.Password `conf:"env"`
}
