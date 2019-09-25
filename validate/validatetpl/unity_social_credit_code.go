package validatetpl

import (
	"regexp"
)

const (
	InvalidUnitySocialCreditCodeType  = "社会信用代码类型错误"
	InvalidUnitySocialCreditCodeValue = "无效的社会信用代码"
)

var (
	unity_social_credit_code_regexp = regexp.MustCompile(`^\w{18}$`)
)

func ValidateUnitySocialCreditCode(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidUnitySocialCreditCodeType
	}
	if !unity_social_credit_code_regexp.MatchString(s) {
		return false, InvalidUnitySocialCreditCodeValue
	}
	return true, ""
}

func ValidateUnitySocialCreditCodeOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidUnitySocialCreditCodeType
	}
	if s != "" && !unity_social_credit_code_regexp.MatchString(s) {
		return false, InvalidUnitySocialCreditCodeValue
	}
	return true, ""
}
