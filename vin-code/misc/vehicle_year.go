package misc

import (
	"strconv"
	"time"
)

func GetModelYear(y rune) int {
	errYear := 0
	v, ok := vinYear[y]
	if !ok {
		return errYear
	}
	curr := time.Now().Year()
	if v[1] <= curr {
		return v[1]
	}
	return v[0]
}

func GetModelYearStr(y rune) string {
	yearNum := GetModelYear(y)
	return strconv.FormatInt(int64(yearNum), 10)
}
