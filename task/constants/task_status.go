package constants

//go:generate libtools gen enum TaskStatus
// swagger:enum
type TaskStatus uint8

// 任务状态
const (
	TASK_STATUS_UNKNOWN     TaskStatus = iota
	TASK_STATUS__INIT        // 就绪
	TASK_STATUS__PENGDING    // 已分发
	TASK_STATUS__PROCESSING  // 执行中
	TASK_STATUS__SUCCESS     // 已完成
	TASK_STATUS__FAIL        // 失败
	TASK_STATUS__ROLLBACK    // 回滚
)
