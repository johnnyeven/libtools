package constants

import "github.com/johnnyeven/libtools/timelib"

type CronTableInfo struct {
	CronTableID string                 `json:"cronTableID"`
	Channel     string                 `json:"channel"`
	Subject     string                 `json:"subject"`
	Spec        string                 `json:"spec"`
	Args        string                 `json:"args"`
	NextTime    timelib.MySQLTimestamp `json:"nextTime"`
	Desc        string                 `json:"desc"`
}
