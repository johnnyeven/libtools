package amount

import (
	"fmt"
	"math"
	"strconv"
)

// FenToYuan from fen to yuan
func FenToYuan(a int64) float64 {
	return Int64ToFloat64(a, 2)
}

// YuanToFen from yun to fen
func YuanToFen(a float64) int64 {
	return Float64ToInt64(a, 2)
}

// convert int64 to float64
func Int64ToFloat64(n int64, decimal int) float64 {
	var negative bool
	if n < 0 {
		negative = true
		n = -n
	}
	k := int64(math.Pow10(decimal))
	head := n / k
	tail := n % k
	format := "%d.%0" + fmt.Sprintf("%d", decimal) + "d"
	str := fmt.Sprintf(format, head, tail)
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(fmt.Sprintf("convert int64[%d] to float64[decimal:%d] failed[err:%s]", n, decimal, err.Error()))
	}
	if negative {
		f = -f
	}
	return f
}

// convert floa64 to int64
func Float64ToInt64(f float64, decimal int) int64 {
	var negative bool
	if f < 0 {
		negative = true
		f = -f
	}
	k := math.Pow10(decimal)
	str := fmt.Sprintf("%.0f", f*k)
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("convert float64[%f] to float64[decimal:%d] failed[err:%s]", f, decimal, err.Error()))
	}
	if negative {
		n = -n
	}
	return n
}

// Round 四舍五入
func Round(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v*pow)+0.5)) / pow
}

/*
利息罚息计算(通过年化利率计算)
本金-capital
年化利率-rate
计息天数-dateDiff
保留位数-decimals
单位-分
*/
func CountInterestByAnnualizedRate(capital, rate float64, dateDiff, decimals int) float64 {
	return Round(float64(capital*float64(dateDiff)*rate/360), decimals)
}
