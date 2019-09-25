package err

import (
	"github.com/johnnyeven/libtools/httplib"
)

var (
	// 400
	TransIncomeAccountNotFoundError = &httplib.GeneralError{
		Msg:            "收益账户不存在",
		Code:           400001,
		Desc:           "收益账户不存在",
		CanBeErrorTalk: false,
	}
	TransFeeGteAmountError = &httplib.GeneralError{
		Msg:            "手续费大于等于实际金额",
		Code:           400002,
		Desc:           "手续费大于等于实际金额",
		CanBeErrorTalk: false,
	}
	TransIncomeAccountSameToSellerError = &httplib.GeneralError{
		Msg:            "收益账户与卖家相同",
		Code:           400003,
		Desc:           "收益账户与卖家相同",
		CanBeErrorTalk: false,
	}
	TransFeeAmoutAccountBindError = &httplib.GeneralError{
		Msg:            "手续费不为0的情况下，收益账户必须存在",
		Code:           400004,
		Desc:           "手续费不为0的情况下，收益账户必须存在",
		CanBeErrorTalk: false,
	}
	TransQueryParamsBindError = &httplib.GeneralError{
		Msg:            "起止时间有误",
		Code:           400015,
		Desc:           "起止时间有误",
		CanBeErrorTalk: true,
	}
)
