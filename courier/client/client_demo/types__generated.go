package client_demo

import (
	golib_tools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
	golib_tools_courier_swagger "github.com/johnnyeven/libtools/courier/swagger"
	golib_tools_courier_transport_http "github.com/johnnyeven/libtools/courier/transport_http"
)

type Data struct {
	//
	Label string `json:"label"`
	//
	Name string `json:"name"`
}

type ErrorField = golib_tools_courier_status_error.ErrorField

type ErrorFields = golib_tools_courier_status_error.ErrorFields

type File = golib_tools_courier_transport_http.File

type GetByID struct {
	//
	ID string `json:"id"`
	//
	Label []string `json:"label,omitempty"`
	//
	Name string `json:"name,omitempty"`
	//
	Status ResourceStatus `json:"status"`
}

type JSONBytes = golib_tools_courier_swagger.JSONBytes

type ResourceStatus = DemoResourceStatus

type StatusError = golib_tools_courier_status_error.StatusError
