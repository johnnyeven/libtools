package redis

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/mq"
)

func TestRedisBackend(t *testing.T) {
	tt := assert.New(t)

	backend := NewRedisBackend(r.GetCache().Pool, "mq-")

	task := mq.NewTask("c", "test", []byte(fmt.Sprintf("data")))

	taskList := make([]*mq.TaskStatus, 0)
	for i := 0; i < 10; i++ {
		taskStatus := task.Success(nil)
		taskList = append(taskList, taskStatus)
		err := backend.FeedBack(taskStatus)
		tt.NoError(err)
	}

	for i := 0; i < 10; i++ {
		taskStatus, err := backend.GetFeedback()
		tt.NoError(err)
		tt.Equal(taskList[i].ID, taskStatus.ID)
		tt.Equal(mq.StatusSuccess, taskStatus.Status)
	}

	{
		err := backend.Cancel("test")
		tt.NoError(err)
	}

	{
		ok, err := backend.IsCancelled("test")
		tt.NoError(err)
		tt.True(ok)
	}

	{
		err := backend.ClearCancellation("test")
		tt.NoError(err)
	}

	{
		ok, err := backend.IsCancelled("test")
		tt.NoError(err)
		tt.False(ok)
	}
}

func TestRedisChannelMgr(t *testing.T) {
	tt := assert.New(t)

	backend := NewRedisBackend(r.GetCache().Pool, "mq-")

	backend.RegisterChannel("test", []string{"One", "Two", "Tree"})

	channelList, err := backend.ListChannel()
	tt.NoError(err)
	tt.Equal([]string{"test"}, channelList)

	subjectList, err := backend.ListSubject("test")
	tt.NoError(err)
	sort.Strings(subjectList)
	tt.Equal([]string{"One", "Tree", "Two"}, subjectList)
}
