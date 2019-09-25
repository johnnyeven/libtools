package slice_to_excel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/slice_to_excel"
	"github.com/johnnyeven/libtools/timelib"
)

type GroupedColumn struct {
	Int    int    `xlsx:"int1"`
	String string `xlsx:"string1"`
}

type Column struct {
	GroupedColumn
	Int    int                    `xlsx:"int2"`
	String string                 `xlsx:"string2"`
	Time   timelib.MySQLTimestamp `xlsx:"time1"`
	Time2  timelib.MySQLTimestamp `xlsx:"time2"`
}

func TestToExcel(t *testing.T) {
	tt := assert.New(t)
	time, _ := timelib.ParseMySQLTimestampFromString("2017-06-01T10:00:00Z")
	time2 := timelib.MySQLTimestampZero
	rows := []Column{
		{
			GroupedColumn: GroupedColumn{
				Int:    11,
				String: "s1",
			},
			Int:    12,
			String: "s2",
			Time:   time,
			Time2:  time2,
		},
	}

	file, err := slice_to_excel.GetExcel("test", rows)
	tt.Nil(err)

	{
		labelRows := []string{"int1", "string1", "int2", "string2", "time1", "time2"}

		cellValues := []string{}

		for _, cell := range file.Sheet["test"].Rows[0].Cells {
			cellValues = append(cellValues, cell.Value)
		}

		tt.Equal(labelRows, cellValues)

	}
	{
		valueRow := []string{"11", "s1", "12", "s2", "2017-06-01T18:00:00+08:00", ""}

		cellValues := []string{}

		for _, cell := range file.Sheet["test"].Rows[1].Cells {
			cellValues = append(cellValues, cell.Value)
		}

		tt.Equal(valueRow, cellValues)
	}
}
