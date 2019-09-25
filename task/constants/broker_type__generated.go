package constants

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidBrokerType = errors.New("invalid BrokerType")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("BrokerType", map[string]string{
		"GEARMAN": "gearman",
		"KAFKA":   "kafka",
		"REDIS":   "redis",
	})
}

func ParseBrokerTypeFromString(s string) (BrokerType, error) {
	switch s {
	case "":
		return BROKER_TYPE_UNKNOWN, nil
	case "GEARMAN":
		return BROKER_TYPE__GEARMAN, nil
	case "KAFKA":
		return BROKER_TYPE__KAFKA, nil
	case "REDIS":
		return BROKER_TYPE__REDIS, nil
	}
	return BROKER_TYPE_UNKNOWN, InvalidBrokerType
}

func ParseBrokerTypeFromLabelString(s string) (BrokerType, error) {
	switch s {
	case "":
		return BROKER_TYPE_UNKNOWN, nil
	case "gearman":
		return BROKER_TYPE__GEARMAN, nil
	case "kafka":
		return BROKER_TYPE__KAFKA, nil
	case "redis":
		return BROKER_TYPE__REDIS, nil
	}
	return BROKER_TYPE_UNKNOWN, InvalidBrokerType
}

func (BrokerType) EnumType() string {
	return "BrokerType"
}

func (BrokerType) Enums() map[int][]string {
	return map[int][]string{
		int(BROKER_TYPE__GEARMAN): {"GEARMAN", "gearman"},
		int(BROKER_TYPE__KAFKA):   {"KAFKA", "kafka"},
		int(BROKER_TYPE__REDIS):   {"REDIS", "redis"},
	}
}
func (v BrokerType) String() string {
	switch v {
	case BROKER_TYPE_UNKNOWN:
		return ""
	case BROKER_TYPE__GEARMAN:
		return "GEARMAN"
	case BROKER_TYPE__KAFKA:
		return "KAFKA"
	case BROKER_TYPE__REDIS:
		return "REDIS"
	}
	return "UNKNOWN"
}

func (v BrokerType) Label() string {
	switch v {
	case BROKER_TYPE_UNKNOWN:
		return ""
	case BROKER_TYPE__GEARMAN:
		return "gearman"
	case BROKER_TYPE__KAFKA:
		return "kafka"
	case BROKER_TYPE__REDIS:
		return "redis"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*BrokerType)(nil)

func (v BrokerType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidBrokerType
	}
	return []byte(str), nil
}

func (v *BrokerType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseBrokerTypeFromString(string(bytes.ToUpper(data)))
	return
}
