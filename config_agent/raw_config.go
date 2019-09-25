package config_agent

import "github.com/johnnyeven/libtools/sqlx/presets"

type RawConfig struct {
	presets.PrimaryID
	// 业务ID
	ConfigurationID uint64 `json:"configurationID,string"`
	// StackID
	StackID uint64 `json:"stackID,string"`
	// Key
	Key string `json:"key"`
	// Value
	Value string `json:"value"`

	presets.OperateTime
	presets.SoftDelete
}
