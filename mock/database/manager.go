package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/pkg/errors"
)

type MockDB struct {
	mockDB *sql.DB
	mock   sqlmock.Sqlmock
	db     *sqlx.DB
}

func (m *MockDB) Init() error {
	var err error
	m.mockDB, m.mock, err = sqlmock.New()
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDB) Open() *sqlx.DB {
	if m.db == nil {
		m.db, _ = sqlx.Open("", "", func(driverName string, dataSourceName string) (db *sql.DB, err error) {
			return m.mockDB, nil
		})
	}

	return m.db
}

func (m *MockDB) GetDB() *sqlx.DB {
	return m.db
}

func (m *MockDB) LoadTestSuite(path string) error {
	queries := make([]Query, 0)
	err := LoadAndParse(path, &queries)
	if err != nil {
		return err
	}

	for _, q := range queries {
		switch q.Type {
		case "begin":
			m.mock.ExpectBegin()
		case "commit":
			m.mock.ExpectCommit()
		case "rollback":
			m.mock.ExpectRollback()
		case "exec":
			exec := m.mock.ExpectExec(q.ExpectedSQLKeyWord)
			if q.WithArgs != nil {
				exec = exec.WithArgs(convertDBValue(q.WithArgs)...)
			}
			if q.ReturnError != "" {
				exec.WillReturnError(errors.Errorf(q.ReturnError))
			} else if q.ReturnResult != nil {
				exec.WillReturnResult(sqlmock.NewResult(q.ReturnResult.LastInsertID, q.ReturnResult.RowsEffected))
			} else {
				return fmt.Errorf("error or result are all nil")
			}
		case "query":
			query := m.mock.ExpectQuery(q.ExpectedSQLKeyWord)
			if q.WithArgs != nil {
				query = query.WithArgs(convertDBValue(q.WithArgs)...)
			}
			if q.ReturnError != "" {
				query.WillReturnError(errors.Errorf(q.ReturnError))
			} else if q.ReturnRows != nil {
				rows := sqlmock.NewRows(q.ReturnRows.Columns)
				for _, r := range q.ReturnRows.Rows {
					values := make([]driver.Value, 0)
					for _, v := range r {
						values = append(values, v)
					}
					rows.AddRow(values...)
				}
				query.WillReturnRows(rows)
			} else {
				return fmt.Errorf("error or rows are all nil")
			}
		default:
			return fmt.Errorf("not supported type %s", q.Type)
		}
	}

	return nil
}

func convertDBValue(data []interface{}) []driver.Value {
	args := make([]driver.Value, 0)
	for _, a := range data {
		var v driver.Value
		f := a.(float64)
		if isValidFloatValue(f) {
			v = f
		} else {
			v = int64(f)
		}
		args = append(args, v)
	}

	return args
}

func isValidFloatValue(v float64) bool {
	compare := v
	if v == float64(int64(compare)) {
		return false
	}

	return true
}
