package httpx

import (
	"net/http"
)

type MethodGet string

func (m MethodGet) Method() string {
	return http.MethodGet
}

type MethodHead string

func (m MethodHead) Method() string {
	return http.MethodHead
}

type MethodPost string

func (m MethodPost) Method() string {
	return http.MethodPost
}

type MethodPut string

func (m MethodPut) Method() string {
	return http.MethodPut
}

type MethodPatch string

func (m MethodPatch) Method() string {
	return http.MethodPatch
}

type MethodDelete string

func (m MethodDelete) Method() string {
	return http.MethodDelete
}
