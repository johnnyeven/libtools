package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"errors"

	golib_tools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidState = errors.New("invalid State")

func init() {
	golib_tools_courier_enumeration.RegisterEnums("State", map[string]string{
		"ACTIVE": "激活",
		"CLOSED": "关闭",
	})
}

func ParseStateFromString(s string) (State, error) {
	switch s {
	case "":
		return STATE_UNKNOWN, nil
	case "ACTIVE":
		return STATE__ACTIVE, nil
	case "CLOSED":
		return STATE__CLOSED, nil
	}
	return STATE_UNKNOWN, InvalidState
}

func ParseStateFromLabelString(s string) (State, error) {
	switch s {
	case "":
		return STATE_UNKNOWN, nil
	case "激活":
		return STATE__ACTIVE, nil
	case "关闭":
		return STATE__CLOSED, nil
	}
	return STATE_UNKNOWN, InvalidState
}

func (State) EnumType() string {
	return "State"
}

func (State) Enums() map[int][]string {
	return map[int][]string{
		int(STATE__ACTIVE): {"ACTIVE", "激活"},
		int(STATE__CLOSED): {"CLOSED", "关闭"},
	}
}
func (v State) String() string {
	switch v {
	case STATE_UNKNOWN:
		return ""
	case STATE__ACTIVE:
		return "ACTIVE"
	case STATE__CLOSED:
		return "CLOSED"
	}
	return "UNKNOWN"
}

func (v State) Label() string {
	switch v {
	case STATE_UNKNOWN:
		return ""
	case STATE__ACTIVE:
		return "激活"
	case STATE__CLOSED:
		return "关闭"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*State)(nil)

func (v State) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidState
	}
	return []byte(str), nil
}

func (v *State) UnmarshalText(data []byte) (err error) {
	*v, err = ParseStateFromString(string(bytes.ToUpper(data)))
	return
}

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*State)(nil)

func (v *State) Scan(src interface{}) error {
	integer, err := golib_tools_courier_enumeration.AsInt64(src, STATE_OFFSET)
	if err != nil {
		return err
	}
	*v = State(integer - STATE_OFFSET)
	return nil
}

func (v State) Value() (driver.Value, error) {
	return int64(v) + STATE_OFFSET, nil
}
