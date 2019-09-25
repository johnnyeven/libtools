package conf

import (
	"fmt"
)

func RancherInternal(stackName string, serviceName string) rancherInternal {
	return rancherInternal{
		StackName:   stackName,
		ServiceName: serviceName,
	}
}

type rancherInternal struct {
	StackName   string
	ServiceName string
}

func (r rancherInternal) String() string {
	return fmt.Sprintf("%s.%s.rancher.internal", r.ServiceName, r.StackName)
}
