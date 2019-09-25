package client_configurations

import (
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
	github_com_johnnyeven_libtools_courier_swagger "github.com/johnnyeven/libtools/courier/swagger"
	github_com_johnnyeven_libtools_sqlx_presets "github.com/johnnyeven/libtools/sqlx/presets"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

type Configuration struct {
	//
	PrimaryID
	//
	OperateTime
	//
	SoftDelete
	// 业务ID
	ConfigurationID uint64 `json:"configurationID,string"`
	// Key
	Key string `json:"key"`
	// StackID
	StackID uint64 `json:"stackID,string"`
	// Value
	Value string `json:"value"`
}

type ConfigurationList []Configuration

type CreateConfigurationBody struct {
	// Key
	Key string `json:"key"`
	// StackID
	StackID uint64 `json:"stackID,string"`
	// Value
	Value string `json:"value"`
}

type ErrorField = github_com_johnnyeven_libtools_courier_status_error.ErrorField

type ErrorFields = github_com_johnnyeven_libtools_courier_status_error.ErrorFields

type GetConfigurationResult struct {
	//
	Data ConfigurationList `json:"data"`
	//
	Total int32 `json:"total"`
}

type JSONBytes = github_com_johnnyeven_libtools_courier_swagger.JSONBytes

type MySQLTimestamp = github_com_johnnyeven_libtools_timelib.MySQLTimestamp

type OperateTime = github_com_johnnyeven_libtools_sqlx_presets.OperateTime

type PrimaryID = github_com_johnnyeven_libtools_sqlx_presets.PrimaryID

type SoftDelete = github_com_johnnyeven_libtools_sqlx_presets.SoftDelete

type StatusError = github_com_johnnyeven_libtools_courier_status_error.StatusError
