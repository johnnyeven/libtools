package timelib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	tt := assert.New(t)

	t0, _ := time.Parse(time.RFC3339, "2017-03-27T23:58:59+08:00")
	dt := MySQLTimestamp(t0)
	tt.Equal("2017-03-27T23:58:59+08:00", dt.String())
	tt.Equal("2017-03-27T23:58:59+08:00", dt.Format(time.RFC3339))
	tt.Equal(int64(1490630339), dt.Unix())
	tt.Equal(MySQLTimestampUnixZero.Unix(), int64(0))
	tt.Equal(MySQLTimestampUnixZero.IsZero(), true)
	tt.Equal("1970-01-01T08:00:00+08:00", MySQLTimestampUnixZero.String())

	{
		dateString, err := dt.MarshalText()
		tt.NoError(err)
		tt.Equal("2017-03-27T23:58:59+08:00", string(dateString))

		dt2 := MySQLTimestampZero
		tt.True(dt2.IsZero())
		err = dt2.UnmarshalText(dateString)
		tt.NoError(err)
		tt.Equal(dt, dt2)
		tt.False(dt2.IsZero())
	}
	{
		value, err := dt.Value()
		tt.NoError(err)
		tt.Equal(int64(1490630339), value.(int64))

		dt2 := MySQLTimestampZero
		tt.True(dt2.IsZero())
		err = dt2.Scan(value)
		tt.NoError(err)
		tt.Equal(dt, dt2)
		tt.False(dt2.IsZero())
	}
	{
		dt3 := MySQLTimestampZero
		err := dt3.UnmarshalText([]byte("\"\""))
		tt.NoError(err)
	}
	{
		dt3 := MySQLTimestampZero
		err := dt3.UnmarshalText([]byte("0"))
		tt.NoError(err)
	}
}
