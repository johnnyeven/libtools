package trans

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"golib/gorm"
)

type Role struct {
	// 主键id
	// Read Only : true
	Id uint64 `gorm:"primary_key;column:F_id" sql:"type:bigint(64) unsigned auto_increment;not null" json:"id"`
	// 角色名称,唯一
	// Maximum length : 64
	Name       string `gorm:"column:F_name" sql:"type:varchar(64);not null;unique_index:I_name" json:"name"`
	CreateTime uint64 `gorm:"column:F_create_time" sql:"type:bigint(64) unsigned;not null;default:0" json:"-"`
	UpdateTime uint64 `gorm:"column:F_update_time" sql:"type:bigint(64) unsigned;not null;default:0;index:I_update_time" json:"-"`
	Enabled    uint8  `gorm:"column:F_enabled" sql:"type:tinyint(8) unsigned;not null;default:1" json:"-"`
}

func (role Role) TableName() string {
	return "access.t_role"
}

func TestTrans(t *testing.T) {
	//init mysql connect pool
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	if err != nil {
		fmt.Printf("Error:%s", err.Error())
		return
	}

	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(5)
	db.SingularTable(true)

	err = db.DB().Ping()
	if err != nil {
		fmt.Printf("Error:%s", err.Error())
		return
	}

	var test_func = func(test_db *gorm.DB) error {
		role := Role{Name: "richard"}
		if err := test_db.Create(&role).Error; err != nil {
			return fmt.Errorf("Error:%s", err.Error())
		}
		panic("test cache transaction")
	}

	err = ExecTransaction(&db, test_func)
	if err != nil {
		fmt.Printf("56:Error:%s", err.Error())
	}
}
