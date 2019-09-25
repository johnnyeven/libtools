package types

// swagger:enum
type State int64

// db -3 = state = 1
const STATE_OFFSET = int64(-3 - STATE__ACTIVE)

const (
	STATE_UNKNOWN State = iota
	STATE__ACTIVE       // 激活
	STATE__CLOSED       // 关闭
)
