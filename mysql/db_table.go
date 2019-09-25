package mysql

import (
	"reflect"

	"github.com/sirupsen/logrus"

	"golib/gorm"

	"github.com/johnnyeven/libtools/env"
)

type Model interface {
	TableName() string
}

func NewDBTable() *DBTable {
	return &DBTable{}
}

type DBTable struct {
	Name string
	list []Model
}

func (dbTable *DBTable) SetName(name string) {
	dbTable.Name = name
}

func (dbTable *DBTable) Register(model Model) {
	rv := reflect.ValueOf(model)
	if rv.Elem().Kind() != reflect.Struct || rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("register model failed")
	}
	dbTable.list = append(dbTable.list, model)
}

func (dbTable *DBTable) AutoMigrate(db *gorm.DB) (err error) {
	goEnv := env.GetRuntimeEnv()

	if goEnv == env.ONLINE || goEnv == env.PRE {
		return nil
	}

	for _, model := range dbTable.list {
		err := db.AutoMigrate(model).Error
		if err != nil {
			logrus.Errorf("%s automigrate error[%s]", model, err.Error())
		}
	}

	return nil
}
