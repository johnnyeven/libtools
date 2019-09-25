package client_demo

import (
	"bytes"
	"encoding"
	"errors"

	golib_tools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

// swagger:enum
type DemoResourceStatus uint

const (
	DEMO_RESOURCE_STATUS_UNKNOWN           DemoResourceStatus = iota
	DEMO_RESOURCE_STATUS__RUNNING                             // 运行中
	DEMO_RESOURCE_STATUS__RECYCLING                           // 回收中
	DEMO_RESOURCE_STATUS__TRANSFORMING_OUT                    // 资源转出中
	DEMO_RESOURCE_STATUS__TRANSFORMING_IN                     // 资源转入中
)

var InvalidDemoResourceStatus = errors.New("invalid DemoResourceStatus")

func init() {
	golib_tools_courier_enumeration.RegisterEnums("DemoResourceStatus", map[string]string{
		"RUNNING":          "运行中",
		"RECYCLING":        "回收中",
		"TRANSFORMING_OUT": "资源转出中",
		"TRANSFORMING_IN":  "资源转入中",
	})
}

func ParseDemoResourceStatusFromString(s string) (DemoResourceStatus, error) {
	switch s {
	case "":
		return DEMO_RESOURCE_STATUS_UNKNOWN, nil
	case "RUNNING":
		return DEMO_RESOURCE_STATUS__RUNNING, nil
	case "RECYCLING":
		return DEMO_RESOURCE_STATUS__RECYCLING, nil
	case "TRANSFORMING_OUT":
		return DEMO_RESOURCE_STATUS__TRANSFORMING_OUT, nil
	case "TRANSFORMING_IN":
		return DEMO_RESOURCE_STATUS__TRANSFORMING_IN, nil
	}
	return DEMO_RESOURCE_STATUS_UNKNOWN, InvalidDemoResourceStatus
}

func ParseDemoResourceStatusFromLabelString(s string) (DemoResourceStatus, error) {
	switch s {
	case "":
		return DEMO_RESOURCE_STATUS_UNKNOWN, nil
	case "运行中":
		return DEMO_RESOURCE_STATUS__RUNNING, nil
	case "回收中":
		return DEMO_RESOURCE_STATUS__RECYCLING, nil
	case "资源转出中":
		return DEMO_RESOURCE_STATUS__TRANSFORMING_OUT, nil
	case "资源转入中":
		return DEMO_RESOURCE_STATUS__TRANSFORMING_IN, nil
	}
	return DEMO_RESOURCE_STATUS_UNKNOWN, InvalidDemoResourceStatus
}

func (DemoResourceStatus) EnumType() string {
	return "DemoResourceStatus"
}

func (DemoResourceStatus) Enums() map[int][]string {
	return map[int][]string{
		int(DEMO_RESOURCE_STATUS__RUNNING):          {"RUNNING", "运行中"},
		int(DEMO_RESOURCE_STATUS__RECYCLING):        {"RECYCLING", "回收中"},
		int(DEMO_RESOURCE_STATUS__TRANSFORMING_OUT): {"TRANSFORMING_OUT", "资源转出中"},
		int(DEMO_RESOURCE_STATUS__TRANSFORMING_IN):  {"TRANSFORMING_IN", "资源转入中"},
	}
}
func (v DemoResourceStatus) String() string {
	switch v {
	case DEMO_RESOURCE_STATUS_UNKNOWN:
		return ""
	case DEMO_RESOURCE_STATUS__RUNNING:
		return "RUNNING"
	case DEMO_RESOURCE_STATUS__RECYCLING:
		return "RECYCLING"
	case DEMO_RESOURCE_STATUS__TRANSFORMING_OUT:
		return "TRANSFORMING_OUT"
	case DEMO_RESOURCE_STATUS__TRANSFORMING_IN:
		return "TRANSFORMING_IN"
	}
	return "UNKNOWN"
}

func (v DemoResourceStatus) Label() string {
	switch v {
	case DEMO_RESOURCE_STATUS_UNKNOWN:
		return ""
	case DEMO_RESOURCE_STATUS__RUNNING:
		return "运行中"
	case DEMO_RESOURCE_STATUS__RECYCLING:
		return "回收中"
	case DEMO_RESOURCE_STATUS__TRANSFORMING_OUT:
		return "资源转出中"
	case DEMO_RESOURCE_STATUS__TRANSFORMING_IN:
		return "资源转入中"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*DemoResourceStatus)(nil)

func (v DemoResourceStatus) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidDemoResourceStatus
	}
	return []byte(str), nil
}

func (v *DemoResourceStatus) UnmarshalText(data []byte) (err error) {
	*v, err = ParseDemoResourceStatusFromString(string(bytes.ToUpper(data)))
	return
}
