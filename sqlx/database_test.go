package sqlx_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/runner"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/johnnyeven/libtools/sqlx/builder"
	"github.com/johnnyeven/libtools/timelib"
)

var db *sqlx.DB

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	db, _ = sqlx.Open("logger:mysql", "root:root@tcp(localhost:3306)/?charset=utf8&parseTime=true&interpolateParams=true&autocommit=true&loc=Local")
}

type TableOperateTime struct {
	CreatedAt timelib.MySQLDatetime `db:"F_created_at" sql:"datetime(6) NOT NULL DEFAULT '0' ON UPDATE CURRENT_TIMESTAMP(6)" `
	UpdatedAt int64                 `db:"F_updated_at" sql:"bigint(64) NOT NULL DEFAULT '0'"`
}

func (t *TableOperateTime) BeforeUpdate() {
	time.Now()
	t.UpdatedAt = time.Now().UnixNano()
}

func (t *TableOperateTime) BeforeInsert() {
	t.CreatedAt = timelib.MySQLDatetime(time.Now())
	t.UpdatedAt = t.CreatedAt.Unix()
}

type Gender int

const (
	GenderMale Gender = iota + 1
	GenderFemale
)

func (Gender) EnumType() string {
	return "Gender"
}

func (Gender) Enums() map[int][]string {
	return map[int][]string{
		int(GenderMale):   {"male", "男"},
		int(GenderFemale): {"female", "女"},
	}
}

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	}
	return ""
}

// @def primary ID
// @def index I_nickname Nickname Name
// @def unique_index I_name Name
type User struct {
	ID       uint64 `db:"F_id" sql:"bigint(64) unsigned NOT NULL AUTO_INCREMENT"`
	Name     string `db:"F_name" sql:"varchar(255) binary NOT NULL DEFAULT ''"`
	Username string `db:"F_username" sql:"varchar(255)"`
	Nickname string `db:"F_nickname" sql:"varchar(255) CHARACTER SET latin1 binary NOT NULL DEFAULT ''"`
	Gender   Gender `db:"F_gender" sql:"int(32) NOT NULL DEFAULT '0'"`

	TableOperateTime
}

func (user *User) Comments() map[string]string {
	return map[string]string{
		"Name": "姓名",
	}
}

func (user *User) AfterInsert(result sql.Result) error {
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(lastInsertID)
	return nil
}

func (user *User) TableName() string {
	return "t_user"
}

func (user *User) PrimaryKey() sqlx.FieldNames {
	return []string{"ID"}
}

func (user *User) Indexes() sqlx.Indexes {
	return sqlx.Indexes{
		"I_nickname": {"Nickname", "Name"},
	}
}

func (user *User) UniqueIndexes() sqlx.Indexes {
	return sqlx.Indexes{
		"I_name": {"Name"},
	}
}

type User2 struct {
	User
	Age int32 `db:"F_age" sql:"int(32) NOT NULL DEFAULT '0'"`
}

func (user2 *User2) Indexes() sqlx.Indexes {
	return sqlx.Indexes{}
}

func TestMigrate(t *testing.T) {
	tt := assert.New(t)

	os.Setenv("PROJECT_FEATURE", "test")
	dbTest := sqlx.NewFeatureDatabase("test_for_migrate")
	defer func() {
		err := db.Do(dbTest.Drop()).Err()
		tt.Nil(err)
	}()

	{
		dbTest.Register(&User{})
		err := dbTest.MigrateTo(db, false)
		tt.NoError(err)
	}
	{
		dbTest.Register(&User{})
		err := dbTest.MigrateTo(db, false)
		tt.NoError(err)
	}
	{
		dbTest.Register(&User2{})
		err := dbTest.MigrateTo(db, false)
		tt.NoError(err)
	}

	{
		dbTest.Register(&User{})
		err := dbTest.MigrateTo(db, false)
		tt.NoError(err)
	}
}

func TestCRUD(t *testing.T) {
	tt := assert.New(t)

	dbTest := sqlx.NewDatabase("test")
	defer func() {
		err := db.Do(dbTest.Drop()).Err()
		tt.Nil(err)
	}()

	userTable := dbTest.Register(&User{})
	err := dbTest.MigrateTo(db, false)
	tt.Nil(err)

	{
		user := User{
			Name:   uuid.New().String(),
			Gender: GenderMale,
		}
		user.BeforeInsert()
		dbRet := db.Do(dbTest.Insert(&user).Comment("InsertUser"))
		err := dbRet.Err()
		tt.Nil(err)
		user.AfterInsert(dbRet.Result)
		tt.NotEmpty(user.ID)

		{
			user.Gender = GenderFemale
			user.BeforeUpdate()
			err := db.Do(
				dbTest.Update(&user).
					Where(
						userTable.F("Name").Eq(user.Name),
					).
					Comment("UpdateUser"),
			).
				Err()
			tt.Nil(err)
		}

		{
			userForSelect := User{}
			err := db.Do(
				userTable.Select().Where(userTable.F("Name").Eq(user.Name)).Comment("FindUser"),
			).
				Scan(&userForSelect).
				Err()
			tt.Nil(err)
			tt.Equal(userForSelect.Name, user.Name)
			tt.Equal(userForSelect.CreatedAt.Unix(), user.CreatedAt.Unix())
		}

		{
			user.BeforeInsert()
			err := db.Do(dbTest.Insert(&user).Comment("Insert Conflict")).Err()
			t.Log(err)
			tt.True(sqlx.DBErr(err).IsConflict())

			{
				err := db.Do(
					dbTest.Insert(&user).
						OnDuplicateKeyUpdate(
							userTable.AssignsByFieldValues(builder.FieldValues{
								"Gender": GenderMale,
							})...,
						).Comment("InsertUserOnDuplicate"),
				).Err()
				tt.Nil(err)
			}
		}
	}

}

func TestSelect(t *testing.T) {
	tt := assert.New(t)

	dbTest := sqlx.NewDatabase("test2")
	defer func() {
		err := db.Do(dbTest.Drop()).Err()
		tt.Nil(err)
	}()

	table := dbTest.Register(&User{})
	err := dbTest.MigrateTo(db, false)
	tt.Nil(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r := runner.NewRunner(ctx, "create data", 30)

	for i := 0; i < 10; i++ {
		r.Add(func(ctx context.Context) error {
			user := User{
				Name:   uuid.New().String(),
				Gender: GenderMale,
			}
			user.BeforeInsert()
			return db.Do(dbTest.Insert(&user).Comment("InsertUser")).Scan(&user).Err()
		}, fmt.Sprintf("%d", i))
	}

	tt.NoError(r.Commit())

	{
		users := make([]User, 0)
		err := db.Do(table.Select().Where(table.F("Gender").Eq(GenderMale))).Scan(&users).Err()
		tt.NoError(err)
		tt.Len(users, 10)
	}
	{
		user := User{}
		err := db.Do(table.Select().Where(table.F("ID").Eq(11))).Scan(&user).Err()
		tt.True(sqlx.DBErr(err).IsNotFound())
	}
	{
		count := 0
		err := db.Do(
			table.Select().For(builder.Count(builder.Star())),
		).Scan(&count).Err()
		tt.NoError(err)
		tt.Equal(10, count)
	}
	{
		user := &User{}
		err := db.Do(table.Select().Where(table.F("Gender").Eq(GenderMale))).Scan(user).Err()
		tt.Error(err)
	}
}
