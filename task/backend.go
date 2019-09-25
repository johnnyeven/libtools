package task

import (
	"github.com/johnnyeven/libtools/task/constants"
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type Backend struct {
	agent *Agent
}

func NewBackend(agent *Agent) *Backend {
	return &Backend{
		agent,
	}
}

func (b *Backend) Feedback(f *constants.TaskFeedback) error {
	bytes, err := json.Marshal(f)
	if err != nil {
		return err
	}

	task := &constants.Task{
		ID:         uuid.New().String(),
		Channel:    b.agent.CentralChannel,
		Subject:    DefaultFeedbackSubject,
		Data:       bytes,
		CreateTime: time.Now(),
	}

	return b.agent.SendTask(task)
}
