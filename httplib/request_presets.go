package httplib

import (
	"github.com/johnnyeven/libtools/timelib"
)

type Pager struct {
	// 分页大小
	// 默认为 10，-1 为查询所有
	Size int32 `name:"size" in:"query" default:"10"  validate:"@int32[-1,50]"`
	// 分页偏移
	// 默认为 0
	Offset int32 `name:"offset,omitempty" in:"query" validate:"@int32[0,]"`
}

type CreateTimeRange struct {
	// 创建起始时间
	CreateStartTime timelib.MySQLTimestamp `name:"createStartTime,omitempty" in:"query"`
	// 创建终止时间
	CreateEndTime timelib.MySQLTimestamp `name:"createEndTime,omitempty" in:"query"`
}

func (createTimeRange CreateTimeRange) ValidateCreateEndTime() string {
	if !createTimeRange.CreateEndTime.IsZero() {
		if createTimeRange.CreateEndTime.Unix() < createTimeRange.CreateStartTime.Unix() {
			return "终止时间不得小于开始时间"
		}
	}
	return ""
}

type UpdateTimeRange struct {
	// 更新起始时间
	UpdateStartTime timelib.MySQLTimestamp `name:"updateStartTime" in:"query" default:"" `
	// 更新终止时间
	UpdateEndTime timelib.MySQLTimestamp `name:"updateEndTime" in:"query" default:""`
}

func (updateTimeRange UpdateTimeRange) ValidateUpdateEndTime() string {
	if !updateTimeRange.UpdateEndTime.IsZero() {
		if updateTimeRange.UpdateEndTime.Unix() < updateTimeRange.UpdateStartTime.Unix() {
			return "终止时间不得小于开始时间"
		}
	}
	return ""
}
