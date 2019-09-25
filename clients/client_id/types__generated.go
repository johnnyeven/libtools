package client_id

import (
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
	github_com_johnnyeven_libtools_courier_swagger "github.com/johnnyeven/libtools/courier/swagger"
)

type ErrorField = github_com_johnnyeven_libtools_courier_status_error.ErrorField

type ErrorFields = github_com_johnnyeven_libtools_courier_status_error.ErrorFields

type JSONBytes = github_com_johnnyeven_libtools_courier_swagger.JSONBytes

type StatusError = github_com_johnnyeven_libtools_courier_status_error.StatusError

type UniqueID struct {
	//
	ID uint64 `json:"id,string"`
}
