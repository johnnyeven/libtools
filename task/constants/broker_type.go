package constants

//go:generate libtools gen enum BrokerType
// swagger:enum
type BrokerType uint8

// Broker类型
const (
	BROKER_TYPE_UNKNOWN  BrokerType = iota
	BROKER_TYPE__GEARMAN  // gearman
	BROKER_TYPE__REDIS    // redis
	BROKER_TYPE__KAFKA    // kafka
)
