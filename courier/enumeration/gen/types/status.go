package types

// swagger:enum
type Status int64

const (
	STATUS_UNKNOWN Status = iota
	STATUS__ACTIVE        // 激活
	STATUS__CLOSED        // 关闭
)
