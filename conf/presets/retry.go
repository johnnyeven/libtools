package presets

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Retry struct {
	Repeats  int
	Interval time.Duration
}

func (r Retry) MarshalDefaults(v interface{}) {
	if retry, ok := v.(*Retry); ok {
		if retry.Repeats == 0 {
			retry.Repeats = 3
		}
		if retry.Interval == 0 {
			retry.Interval = 10 * time.Second
		}
	}
}

func (r Retry) Do(exec func() error) (err error) {
	if r.Repeats <= 0 {
		err = exec()
		return
	}
	for i := 0; i < r.Repeats; i++ {
		err = exec()
		if err != nil {
			logrus.Warningf("retry in seconds [%d]", r.Interval)
			time.Sleep(r.Interval)
		}
	}
	return
}
