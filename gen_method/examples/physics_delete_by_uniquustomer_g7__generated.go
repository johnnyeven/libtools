package examples

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"golib/gorm"

	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/libtools/duration"
	"github.com/johnnyeven/libtools/mysql/dberr"
	"github.com/johnnyeven/libtools/timelib"
)

type PhysicsDeleteByUniquustomerG7List []PhysicsDeleteByUniquustomerG7

func init() {
	DBTable.Register(&PhysicsDeleteByUniquustomerG7{})
}

func (pdbug PhysicsDeleteByUniquustomerG7) TableName() string {
	table_name := "t_physics_delete_by_uniquustomer_g7"
	if DBTable.Name == "" {
		return table_name
	}
	return DBTable.Name + "." + table_name
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) BatchFetchByCreateTimeList(db *gorm.DB, createTimeList []timelib.MySQLTimestamp) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.BatchFetchByCreateTimeList",
	})()

	if len(createTimeList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_create_time in (?) and F_enabled = ?", createTimeList, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) BatchFetchByCustomerIDList(db *gorm.DB, customerIDList []uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.BatchFetchByCustomerIDList",
	})()

	if len(customerIDList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_customer_id in (?) and F_enabled = ?", customerIDList, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) BatchFetchByG7sOrgCodeList(db *gorm.DB, g7sOrgCodeList []string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.BatchFetchByG7sOrgCodeList",
	})()

	if len(g7sOrgCodeList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_g7s_org_code in (?) and F_enabled = ?", g7sOrgCodeList, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) BatchFetchByG7sUserIDList(db *gorm.DB, g7sUserIDList []string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.BatchFetchByG7sUserIDList",
	})()

	if len(g7sUserIDList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_g7s_user_id in (?) and F_enabled = ?", g7sUserIDList, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) BatchFetchByUpdateTimeList(db *gorm.DB, updateTimeList []timelib.MySQLTimestamp) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.BatchFetchByUpdateTimeList",
	})()

	if len(updateTimeList) == 0 {
		return nil
	}

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_update_time in (?) and F_enabled = ?", updateTimeList, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7) Create(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.Create",
	})()

	if time.Time(pdbug.CreateTime).IsZero() {
		pdbug.CreateTime = timelib.MySQLTimestamp(time.Now())
	}

	if time.Time(pdbug.UpdateTime).IsZero() {
		pdbug.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	pdbug.Enabled = enumeration.Bool(enumeration.BOOL__TRUE)
	err := db.Table(pdbug.TableName()).Create(pdbug).Error
	if err != nil {
		if mysql_err, ok := err.(*mysql.MySQLError); !ok {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordCreateFailedError
		} else if mysql_err.Number != dberr.DuplicateEntryErrNumber {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordCreateFailedError
		} else {
			return dberr.RecordConflictError
		}
	} else {
		return nil
	}
}

type PhysicsDeleteByUniquustomerG7DBFieldData struct {
	G7sUserID  string
	CustomerID string
	G7sOrgCode string
	CreateTime string
	UpdateTime string
	Enabled    string
}

// FetchNoneUniqueIndexFields without Enabled and CreateTime field.
func (pdbugdbfd *PhysicsDeleteByUniquustomerG7DBFieldData) FetchNoneUniqueIndexFields() []string {
	return []string{
		"F_g7s_org_code", "F_update_time",
	}
}

func (pdbug PhysicsDeleteByUniquustomerG7) DBField() *PhysicsDeleteByUniquustomerG7DBFieldData {
	return &PhysicsDeleteByUniquustomerG7DBFieldData{
		G7sUserID:  "F_g7s_user_id",
		CustomerID: "F_customer_id",
		G7sOrgCode: "F_g7s_org_code",
		CreateTime: "F_create_time",
		UpdateTime: "F_update_time",
		Enabled:    "F_enabled",
	}
}

var PhysicsDeleteByUniquustomerG7StructFieldAndDBFieldRelate = map[string]string{
	"G7sUserID":  "F_g7s_user_id",
	"CustomerID": "F_customer_id",
	"G7sOrgCode": "F_g7s_org_code",
	"CreateTime": "F_create_time",
	"UpdateTime": "F_update_time",
	"Enabled":    "F_enabled",
}

var PhysicsDeleteByUniquustomerG7DBFieldAndStructFieldRelate = map[string]string{
	"F_g7s_user_id":  "G7sUserID",
	"F_customer_id":  "CustomerID",
	"F_g7s_org_code": "G7sOrgCode",
	"F_create_time":  "CreateTime",
	"F_update_time":  "UpdateTime",
	"F_enabled":      "Enabled",
}

// CreateOnDuplicateWithUpdateFields only update the no unique index field, it return error if updateFields contain unique index field.
// It doesn't update the Enabled and CreateTime field.
func (pdbug *PhysicsDeleteByUniquustomerG7) CreateOnDuplicateWithUpdateFields(db *gorm.DB, updateFields []string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.CreateOnDuplicateWithUpdateFields",
	})()
	if len(updateFields) == 0 {
		return fmt.Errorf("Must have update fields.")
	}

	noUniqueIndexFields := (&PhysicsDeleteByUniquustomerG7DBFieldData{}).FetchNoneUniqueIndexFields()
	if len(noUniqueIndexFields) == 0 {
		return fmt.Errorf("There are no unique fields.")
	}

	var noUniqueIndexFieldsMap = make(map[string]string)
	for _, field := range noUniqueIndexFields {
		noUniqueIndexFieldsMap[field] = ""
	}

	var updateFieldsMap = make(map[string]string)
	for _, field := range updateFields {
		// have unique field in updateFields
		if _, ok := noUniqueIndexFieldsMap[field]; !ok {
			return fmt.Errorf("Field[%s] is unique index or wrong field or Enable field", PhysicsDeleteByUniquustomerG7DBFieldAndStructFieldRelate[field])
		}
		updateFieldsMap[field] = ""
	}

	if time.Time(pdbug.CreateTime).IsZero() {
		pdbug.CreateTime = timelib.MySQLTimestamp(time.Now())
	}

	if time.Time(pdbug.UpdateTime).IsZero() {
		pdbug.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	pdbug.Enabled = enumeration.Bool(enumeration.BOOL__TRUE)

	structType := reflect.TypeOf(pdbug).Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("Instance not struct type.")
	}
	structVal := reflect.ValueOf(pdbug).Elem()

	var param_list []interface{}
	var str_list = []string{"insert into"}
	var insertFieldsStr = pdbug.TableName() + "("
	var placeHolder = "values("
	for i := 0; i < structType.NumField(); i++ {
		if i == 0 {
			insertFieldsStr += PhysicsDeleteByUniquustomerG7StructFieldAndDBFieldRelate[structType.Field(i).Name]
			placeHolder += fmt.Sprintf("%s", "?")
		} else {
			insertFieldsStr += fmt.Sprintf(",%s", PhysicsDeleteByUniquustomerG7StructFieldAndDBFieldRelate[structType.Field(i).Name])
			placeHolder += fmt.Sprintf("%s", ", ?")
		}
		param_list = append(param_list, structVal.Field(i).Interface())
	}
	insertFieldsStr += ")"
	placeHolder += ")"
	str_list = append(str_list, []string{insertFieldsStr, placeHolder, "on duplicate key update"}...)

	var updateStr []string
	for i := 0; i < structType.NumField(); i++ {
		if dbField, ok := PhysicsDeleteByUniquustomerG7StructFieldAndDBFieldRelate[structType.Field(i).Name]; !ok {
			return fmt.Errorf("Wrong field of struct, may be changed field but not regenerate code.")
		} else {
			if _, ok := updateFieldsMap[dbField]; ok {
				updateStr = append(updateStr, fmt.Sprintf("%s = ?", dbField))
				param_list = append(param_list, structVal.Field(i).Interface())
			}
		}
	}
	str_list = append(str_list, strings.Join(updateStr, ","))
	sql := strings.Join(str_list, " ")
	err := db.Exec(sql, param_list...).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordCreateFailedError
	}

	return nil
}

func (pdbug *PhysicsDeleteByUniquustomerG7) DeleteByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.DeleteByG7sUserIDAndCustomerID",
	})()

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Delete(pdbug).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordDeleteFailedError
	} else {
		return nil
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) FetchByCreateTime(db *gorm.DB, createTime timelib.MySQLTimestamp) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByCreateTime",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_create_time = ? and F_enabled = ?", createTime, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) FetchByCustomerID(db *gorm.DB, customerID uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByCustomerID",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_customer_id = ? and F_enabled = ?", customerID, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) FetchByG7sOrgCode(db *gorm.DB, g7sOrgCode string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByG7sOrgCode",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_g7s_org_code = ? and F_enabled = ?", g7sOrgCode, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) FetchByG7sUserID(db *gorm.DB, g7sUserID string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByG7sUserID",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_g7s_user_id = ? and F_enabled = ?", g7sUserID, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7) FetchByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByG7sUserIDAndCustomerID",
	})()

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Find(pdbug).Error
	if err == nil {
		return nil
	} else {
		if err == gorm.RecordNotFound {
			return dberr.RecordNotFoundError
		} else {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordFetchFailedError
		}
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7) FetchByG7sUserIDAndCustomerIDForUpdate(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByG7sUserIDAndCustomerIDForUpdate",
	})()

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Set("gorm:query_option", "FOR UPDATE").Find(pdbug).Error
	if err == nil {
		return nil
	} else {
		if err == gorm.RecordNotFound {
			return dberr.RecordNotFoundError
		} else {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordFetchFailedError
		}
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) FetchByUpdateTime(db *gorm.DB, updateTime timelib.MySQLTimestamp) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchByUpdateTime",
	})()

	err := db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where("F_update_time = ? and F_enabled = ?", updateTime, enumeration.BOOL__TRUE).Find(pdbugl).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (pdbugl *PhysicsDeleteByUniquustomerG7List) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.FetchList",
	})()

	var count int32
	if len(query) == 0 {
		query = append(query, map[string]interface{}{"F_enabled": enumeration.BOOL__TRUE})
	} else {
		if _, ok := query[0]["F_enabled"]; !ok {
			query[0]["F_enabled"] = enumeration.BOOL__TRUE
		}
	}

	if size <= 0 {
		size = -1
		offset = -1
	}
	var err error

	err = db.Table(PhysicsDeleteByUniquustomerG7{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("F_create_time desc").Find(pdbugl).Error

	if err != nil {
		logrus.Errorf("%s", err.Error())
		return 0, dberr.RecordFetchFailedError
	} else {
		return int32(count), nil
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7) SoftDeleteByG7sUserIDAndCustomerID(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.SoftDeleteByG7sUserIDAndCustomerID",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = enumeration.BOOL__FALSE

	if time.Time(pdbug.UpdateTime).IsZero() {
		pdbug.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	err := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Updates(updateMap).Error
	if err != nil {
		if mysql_err, ok := err.(*mysql.MySQLError); !ok {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordDeleteFailedError
		} else if mysql_err.Number != dberr.DuplicateEntryErrNumber {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordDeleteFailedError
		} else {
			logrus.Warningf("%s", err.Error())
			// 物理删除被软删除的数据
			delErr := db.Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Delete(&PhysicsDeleteByUniquustomerG7{}).Error
			if delErr != nil {
				logrus.Errorf("%s", delErr.Error())
				return dberr.RecordDeleteFailedError
			}

			return nil
		}
	} else {
		return nil
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7) UpdateByG7sUserIDAndCustomerIDWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.UpdateByG7sUserIDAndCustomerIDWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = timelib.MySQLTimestamp(time.Now())

	}
	dbRet := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Updates(updateMap)
	err := dbRet.Error
	if err != nil {
		if mysql_err, ok := err.(*mysql.MySQLError); !ok {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordUpdateFailedError
		} else if mysql_err.Number != dberr.DuplicateEntryErrNumber {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordUpdateFailedError
		} else {
			return dberr.RecordConflictError
		}
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Find(&PhysicsDeleteByUniquustomerG7{}).Error
			if findErr == gorm.RecordNotFound {
				return dberr.RecordNotFoundError
			} else if findErr != nil {
				return dberr.RecordUpdateFailedError
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}

func (pdbug *PhysicsDeleteByUniquustomerG7) UpdateByG7sUserIDAndCustomerIDWithStruct(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]PhysicsDeleteByUniquustomerG7.UpdateByG7sUserIDAndCustomerIDWithStruct",
	})()

	if time.Time(pdbug.UpdateTime).IsZero() {
		pdbug.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	dbRet := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Updates(pdbug)
	err := dbRet.Error
	if err != nil {
		if mysql_err, ok := err.(*mysql.MySQLError); !ok {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordUpdateFailedError
		} else if mysql_err.Number != dberr.DuplicateEntryErrNumber {
			logrus.Errorf("%s", err.Error())
			return dberr.RecordUpdateFailedError
		} else {
			return dberr.RecordConflictError
		}
	} else {
		if dbRet.RowsAffected == 0 {
			findErr := db.Table(pdbug.TableName()).Where("F_g7s_user_id = ? and F_customer_id = ? and F_enabled = ?", pdbug.G7sUserID, pdbug.CustomerID, enumeration.BOOL__TRUE).Find(&PhysicsDeleteByUniquustomerG7{}).Error
			if findErr == gorm.RecordNotFound {
				return dberr.RecordNotFoundError
			} else if findErr != nil {
				return dberr.RecordUpdateFailedError
			}
			//存在有效数据记录，返回成功
			return nil
		} else {
			return nil
		}
	}
}
