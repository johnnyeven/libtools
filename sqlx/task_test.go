package sqlx_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/sqlx"
)

func TestWithTasks(t *testing.T) {
	tt := assert.New(t)

	dbTest := sqlx.NewDatabase("test_for_user")
	defer func() {
		err := db.Do(dbTest.Drop()).Err()
		tt.NoError(err)
	}()

	{
		dbTest.Register(&User{})
		err := dbTest.MigrateTo(db, false)
		tt.NoError(err)
	}

	{
		taskList := sqlx.NewTasks(db)

		taskList = taskList.With(func(db *sqlx.DB) error {
			user := User{
				Name:   uuid.New().String(),
				Gender: GenderMale,
			}
			return db.Do(dbTest.Insert(&user).Comment("InsertUser")).Err()
		})

		taskList = taskList.With(func(db *sqlx.DB) error {
			subTaskList := sqlx.NewTasks(db)

			subTaskList = subTaskList.With(func(db *sqlx.DB) error {
				user := User{
					Name:   uuid.New().String(),
					Gender: GenderMale,
				}
				return db.Do(dbTest.Insert(&user).Comment("InsertUser")).Err()
			})

			subTaskList = subTaskList.With(func(db *sqlx.DB) error {
				return fmt.Errorf("rollback")
			})

			return subTaskList.Do()
		})

		err := taskList.Do()
		tt.NotNil(err)
	}

	taskList := sqlx.NewTasks(db)

	taskList = taskList.With(func(db *sqlx.DB) error {
		user := User{
			Name:   uuid.New().String(),
			Gender: GenderMale,
		}
		return db.Do(dbTest.Insert(&user).Comment("InsertUser")).Err()
	})

	taskList = taskList.With(func(db *sqlx.DB) error {
		subTaskList := sqlx.NewTasks(db)

		subTaskList = subTaskList.With(func(db *sqlx.DB) error {
			user := User{
				Name:   uuid.New().String(),
				Gender: GenderMale,
			}
			return db.Do(dbTest.Insert(&user).Comment("InsertUser")).Err()
		})

		subTaskList = subTaskList.With(func(db *sqlx.DB) error {
			user := User{
				Name:   uuid.New().String(),
				Gender: GenderMale,
			}
			return db.Do(dbTest.Insert(&user).Comment("InsertUser")).Err()
		})

		return subTaskList.Do()
	})

	err := taskList.Do()
	tt.NoError(err)
}
