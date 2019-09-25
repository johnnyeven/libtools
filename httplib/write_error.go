package httplib

import (
	"os"

	"github.com/johnnyeven/libtools/courier/status_error"
)

var serviceName = ""

func SetServiceName(name string) {
	serviceName = name
}

func getServiceName() string {
	if serviceName == "" {
		SetServiceName(os.Getenv("PROJECT_NAME"))
	}
	return serviceName
}

func errorWithSource(err *status_error.StatusError) *status_error.StatusError {
	return err.WithSource(getServiceName())
}
