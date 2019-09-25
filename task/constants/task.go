package constants

import (
	"time"
	"fmt"
)

type Task struct {
	ID         string    `json:"id"`
	Channel    string    `json:"channel"`
	Subject    string    `json:"subject"`
	Data       []byte    `json:"data,omitempty"`
	CreateTime time.Time `json:"createTime"`
}

func (t *Task) Init() *TaskFeedback {
	return &TaskFeedback{
		ID:         t.ID,
		Channel:    t.Channel,
		Subject:    t.Subject,
		Data:       t.Data,
		Status:     TASK_STATUS__INIT,
	}
}

func (t *Task) Pending() *TaskFeedback {
	return &TaskFeedback{
		ID:         t.ID,
		Channel:    t.Channel,
		Subject:    t.Subject,
		Status:     TASK_STATUS__PENGDING,
		UpdateTime: time.Now(),
	}
}

func (t *Task) Processing() *TaskFeedback {
	return &TaskFeedback{
		ID:         t.ID,
		Channel:    t.Channel,
		Subject:    t.Subject,
		Status:     TASK_STATUS__PROCESSING,
		UpdateTime: time.Now(),
	}
}

func (t *Task) Success(result interface{}) *TaskFeedback {
	return &TaskFeedback{
		ID:         t.ID,
		Channel:    t.Channel,
		Subject:    t.Subject,
		Status:     TASK_STATUS__SUCCESS,
		UpdateTime: time.Now(),
		Result:     []byte(fmt.Sprintf("%+v", result)),
	}
}

func (t *Task) Fail(err error) *TaskFeedback {
	return &TaskFeedback{
		ID:         t.ID,
		Channel:    t.Channel,
		Subject:    t.Subject,
		Status:     TASK_STATUS__FAIL,
		UpdateTime: time.Now(),
		ErrorTrace: []byte(fmt.Sprintf("%v", err)),
	}
}

func (t *Task) Rollback(err error) *TaskFeedback {
	return &TaskFeedback{
		ID:         t.ID,
		Channel:    t.Channel,
		Subject:    t.Subject,
		Status:     TASK_STATUS__ROLLBACK,
		UpdateTime: time.Now(),
		ErrorTrace: []byte(fmt.Sprintf("%v", err)),
	}
}

type TaskFeedback struct {
	ID         string     `json:"id"`
	Channel    string     `json:"channel"`
	Subject    string     `json:"subject"`
	Data       []byte     `json:"data,omitempty"`
	Status     TaskStatus `json:"status"`
	UpdateTime time.Time  `json:"updateTime"`
	Result     []byte     `json:"result,omitempty"`
	ErrorTrace []byte     `json:"errorTrace,omitempty"`
}
