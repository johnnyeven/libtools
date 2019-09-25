package types

import (
	"bytes"
	"encoding"
	"errors"

	golib_tools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidStatus = errors.New("invalid Status")

func init() {
	golib_tools_courier_enumeration.RegisterEnums("Status", map[string]string{
		"ACTIVE": "激活",
		"CLOSED": "关闭",
	})
}

func ParseStatusFromString(s string) (Status, error) {
	switch s {
	case "":
		return STATUS_UNKNOWN, nil
	case "ACTIVE":
		return STATUS__ACTIVE, nil
	case "CLOSED":
		return STATUS__CLOSED, nil
	}
	return STATUS_UNKNOWN, InvalidStatus
}

func ParseStatusFromLabelString(s string) (Status, error) {
	switch s {
	case "":
		return STATUS_UNKNOWN, nil
	case "激活":
		return STATUS__ACTIVE, nil
	case "关闭":
		return STATUS__CLOSED, nil
	}
	return STATUS_UNKNOWN, InvalidStatus
}

func (Status) EnumType() string {
	return "Status"
}

func (Status) Enums() map[int][]string {
	return map[int][]string{
		int(STATUS__ACTIVE): {"ACTIVE", "激活"},
		int(STATUS__CLOSED): {"CLOSED", "关闭"},
	}
}
func (v Status) String() string {
	switch v {
	case STATUS_UNKNOWN:
		return ""
	case STATUS__ACTIVE:
		return "ACTIVE"
	case STATUS__CLOSED:
		return "CLOSED"
	}
	return "UNKNOWN"
}

func (v Status) Label() string {
	switch v {
	case STATUS_UNKNOWN:
		return ""
	case STATUS__ACTIVE:
		return "激活"
	case STATUS__CLOSED:
		return "关闭"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*Status)(nil)

func (v Status) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidStatus
	}
	return []byte(str), nil
}

func (v *Status) UnmarshalText(data []byte) (err error) {
	*v, err = ParseStatusFromString(string(bytes.ToUpper(data)))
	return
}
