package mq

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewTask(channel string, subject string, data []byte) *Task {
	return &Task{
		MsgHead: MsgHead{
			Channel: channel,
			Subject: subject,
			ID:      genUUID(),
		},
		Data:      data,
		CreatedAt: time.Now(),
	}
}

type Task struct {
	MsgHead
	Data      []byte    `json:"data,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (task *Task) Pending() *TaskStatus {
	return &TaskStatus{
		MsgHead:   task.MsgHead,
		Args:      task.Data,
		UpdatedAt: task.CreatedAt,
		Status:    StatusPending,
	}
}

func (task *Task) Cancelled() *TaskStatus {
	return &TaskStatus{
		MsgHead:   task.MsgHead,
		UpdatedAt: time.Now(),
		Status:    StatusCancelled,
	}
}

func (task *Task) Running() *TaskStatus {
	return &TaskStatus{
		MsgHead:   task.MsgHead,
		UpdatedAt: time.Now(),
		Status:    StatusRunning,
	}
}

func (task *Task) Success(result interface{}) *TaskStatus {
	return &TaskStatus{
		MsgHead:   task.MsgHead,
		UpdatedAt: time.Now(),
		Status:    StatusSuccess,
		Result:    []byte(fmt.Sprintf("%v", result)),
	}
}

func (task *Task) Failed(traceback error) *TaskStatus {
	return &TaskStatus{
		MsgHead:   task.MsgHead,
		UpdatedAt: time.Now(),
		Status:    StatusFailed,
		Traceback: []byte(traceback.Error()),
	}
}

type TaskStatus struct {
	MsgHead
	Status    Status    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	Args      []byte    `json:"args,omitempty"`
	Result    []byte    `json:"result,omitempty"`
	Traceback []byte    `json:"traceback,omitempty"`
}

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCancelled Status = "CANCELED"
	StatusSuccess   Status = "SUCCESS"
	StatusFailed    Status = "FAILED"
)

func genUUID() string {
	return uuid.New().String()
}
