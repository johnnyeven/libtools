package vincode

import (
	"fmt"
)

// VIN 错误定义
type VINError struct {
	ErrCode int
	ErrMsg  string
}

var (
	VINCodeLengthError           = VINError{ErrCode: 20000, ErrMsg: "VIN 长度错误"}
	VINCodeCharacterError        = VINError{ErrCode: 20001, ErrMsg: "VIN 字符非法"}
	VINCodeCheckDigitError       = VINError{ErrCode: 20002, ErrMsg: "VIN 校验失败"}
	VINCodeParseCountryError     = VINError{ErrCode: 20003, ErrMsg: "VIN 解析制造国错误"}
	VINCodeParseManufError       = VINError{ErrCode: 20004, ErrMsg: "VIN 解析制造商错误"}
	VINCodeParseVehicleTypeError = VINError{ErrCode: 20005, ErrMsg: "VIN 解析车用类型错误"}
	VINCodeParseRestriantError   = VINError{ErrCode: 20006, ErrMsg: "VIN 解析约束系统错误"}
	VINCodeParseCarSeriesError   = VINError{ErrCode: 20007, ErrMsg: "VIN 解析车系错误"}
	VINCodeParseDoorTypeError    = VINError{ErrCode: 20008, ErrMsg: "VIN 解析车门错误"}
	VINCodeParseEngineError      = VINError{ErrCode: 20009, ErrMsg: "VIN 解析发动机错误"}
	VINCodeParseAssemblyError    = VINError{ErrCode: 20010, ErrMsg: "VIN 解析装配厂错误"}
	VINCodeParseYearError        = VINError{ErrCode: 20011, ErrMsg: "VIN 解析年份错误"}
	VINCodeParseEmissionError    = VINError{ErrCode: 20012, ErrMsg: "VIN 解析排放错误"}
)

func (e VINError) Error() string {
	return fmt.Sprintf("error code:%d error msg:%s", e.ErrCode, e.ErrMsg)
}
