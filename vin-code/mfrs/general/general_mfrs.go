package general

import (
	"fmt"

	"github.com/johnnyeven/libtools/vin-code/vinrune"
)

// 缺省的解析方式

var (
// 制造国 pos 1
// 制造商 pos 2
// 车用类型 pos 3
// 年份 pos 10
// 顺序号 pos 12~17
)

func GetVISRune(vin string) vinrune.VISRune {
	ret := vinrune.VISRune{}
	vinStr := vin[9:]
	fmt.Sscanf(vinStr, "%c%c%6s", &ret.YearRune, &ret.AssemblyRune, &ret.SequenceNO)
	return ret
}
