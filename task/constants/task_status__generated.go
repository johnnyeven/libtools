package constants

import (
	"bytes"
	"encoding"
	"errors"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
)

var InvalidTaskStatus = errors.New("invalid TaskStatus")

func init() {
	github_com_johnnyeven_libtools_courier_enumeration.RegisterEnums("TaskStatus", map[string]string{
		"FAIL":       "失败",
		"INIT":       "就绪",
		"PENGDING":   "已分发",
		"PROCESSING": "执行中",
		"ROLLBACK":   "回滚",
		"SUCCESS":    "已完成",
	})
}

func ParseTaskStatusFromString(s string) (TaskStatus, error) {
	switch s {
	case "":
		return TASK_STATUS_UNKNOWN, nil
	case "FAIL":
		return TASK_STATUS__FAIL, nil
	case "INIT":
		return TASK_STATUS__INIT, nil
	case "PENGDING":
		return TASK_STATUS__PENGDING, nil
	case "PROCESSING":
		return TASK_STATUS__PROCESSING, nil
	case "ROLLBACK":
		return TASK_STATUS__ROLLBACK, nil
	case "SUCCESS":
		return TASK_STATUS__SUCCESS, nil
	}
	return TASK_STATUS_UNKNOWN, InvalidTaskStatus
}

func ParseTaskStatusFromLabelString(s string) (TaskStatus, error) {
	switch s {
	case "":
		return TASK_STATUS_UNKNOWN, nil
	case "失败":
		return TASK_STATUS__FAIL, nil
	case "就绪":
		return TASK_STATUS__INIT, nil
	case "已分发":
		return TASK_STATUS__PENGDING, nil
	case "执行中":
		return TASK_STATUS__PROCESSING, nil
	case "回滚":
		return TASK_STATUS__ROLLBACK, nil
	case "已完成":
		return TASK_STATUS__SUCCESS, nil
	}
	return TASK_STATUS_UNKNOWN, InvalidTaskStatus
}

func (TaskStatus) EnumType() string {
	return "TaskStatus"
}

func (TaskStatus) Enums() map[int][]string {
	return map[int][]string{
		int(TASK_STATUS__FAIL):       {"FAIL", "失败"},
		int(TASK_STATUS__INIT):       {"INIT", "就绪"},
		int(TASK_STATUS__PENGDING):   {"PENGDING", "已分发"},
		int(TASK_STATUS__PROCESSING): {"PROCESSING", "执行中"},
		int(TASK_STATUS__ROLLBACK):   {"ROLLBACK", "回滚"},
		int(TASK_STATUS__SUCCESS):    {"SUCCESS", "已完成"},
	}
}
func (v TaskStatus) String() string {
	switch v {
	case TASK_STATUS_UNKNOWN:
		return ""
	case TASK_STATUS__FAIL:
		return "FAIL"
	case TASK_STATUS__INIT:
		return "INIT"
	case TASK_STATUS__PENGDING:
		return "PENGDING"
	case TASK_STATUS__PROCESSING:
		return "PROCESSING"
	case TASK_STATUS__ROLLBACK:
		return "ROLLBACK"
	case TASK_STATUS__SUCCESS:
		return "SUCCESS"
	}
	return "UNKNOWN"
}

func (v TaskStatus) Label() string {
	switch v {
	case TASK_STATUS_UNKNOWN:
		return ""
	case TASK_STATUS__FAIL:
		return "失败"
	case TASK_STATUS__INIT:
		return "就绪"
	case TASK_STATUS__PENGDING:
		return "已分发"
	case TASK_STATUS__PROCESSING:
		return "执行中"
	case TASK_STATUS__ROLLBACK:
		return "回滚"
	case TASK_STATUS__SUCCESS:
		return "已完成"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*TaskStatus)(nil)

func (v TaskStatus) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidTaskStatus
	}
	return []byte(str), nil
}

func (v *TaskStatus) UnmarshalText(data []byte) (err error) {
	*v, err = ParseTaskStatusFromString(string(bytes.ToUpper(data)))
	return
}
