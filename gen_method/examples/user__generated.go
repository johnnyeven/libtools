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

type UserList []User

func init() {
	DBTable.Register(&User{})
}

func (ul *UserList) BatchFetchByCityIDList(db *gorm.DB, cityIDList []uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.BatchFetchByCityIDList",
	})()

	if len(cityIDList) == 0 {
		return nil
	}

	err := db.Table(User{}.TableName()).Where("F_city_id in (?) and F_enabled = ?", cityIDList, enumeration.BOOL__TRUE).Find(ul).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (ul *UserList) BatchFetchByIdList(db *gorm.DB, idList []uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.BatchFetchByIdList",
	})()

	if len(idList) == 0 {
		return nil
	}

	err := db.Table(User{}.TableName()).Where("F_id in (?) and F_enabled = ?", idList, enumeration.BOOL__TRUE).Find(ul).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (ul *UserList) BatchFetchByPhoneList(db *gorm.DB, phoneList []string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.BatchFetchByPhoneList",
	})()

	if len(phoneList) == 0 {
		return nil
	}

	err := db.Table(User{}.TableName()).Where("F_phone in (?) and F_enabled = ?", phoneList, enumeration.BOOL__TRUE).Find(ul).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (ul *UserList) BatchFetchByUserIDList(db *gorm.DB, userIDList []uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.BatchFetchByUserIDList",
	})()

	if len(userIDList) == 0 {
		return nil
	}

	err := db.Table(User{}.TableName()).Where("F_user_id in (?) and F_enabled = ?", userIDList, enumeration.BOOL__TRUE).Find(ul).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	} else {
		return nil
	}
}

func (u *User) Create(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.Create",
	})()

	if u.CreateTime.IsZero() {
		u.CreateTime = time.Now()
	}

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	u.Enabled = uint8(enumeration.BOOL__TRUE)
	err := db.Table(u.TableName()).Create(u).Error
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

type UserDBFieldData struct {
	Id         string
	Phone      string
	UserID     string
	Name       string
	CityID     string
	AreaID     string
	Enabled    string
	CreateTime string
	UpdateTime string
}

// FetchNoneUniqueIndexFields without Enabled and CreateTime field.
func (udbfd *UserDBFieldData) FetchNoneUniqueIndexFields() []string {
	return []string{
		"F_city_id", "F_area_id", "F_update_time",
	}
}

func (u User) DBField() *UserDBFieldData {
	return &UserDBFieldData{
		Id:         "F_id",
		Phone:      "F_phone",
		UserID:     "F_user_id",
		Name:       "F_name",
		CityID:     "F_city_id",
		AreaID:     "F_area_id",
		Enabled:    "F_enabled",
		CreateTime: "F_create_time",
		UpdateTime: "F_update_time",
	}
}

var UserStructFieldAndDBFieldRelate = map[string]string{
	"Id":         "F_id",
	"Phone":      "F_phone",
	"UserID":     "F_user_id",
	"Name":       "F_name",
	"CityID":     "F_city_id",
	"AreaID":     "F_area_id",
	"Enabled":    "F_enabled",
	"CreateTime": "F_create_time",
	"UpdateTime": "F_update_time",
}

var UserDBFieldAndStructFieldRelate = map[string]string{
	"F_id":          "Id",
	"F_phone":       "Phone",
	"F_user_id":     "UserID",
	"F_name":        "Name",
	"F_city_id":     "CityID",
	"F_area_id":     "AreaID",
	"F_enabled":     "Enabled",
	"F_create_time": "CreateTime",
	"F_update_time": "UpdateTime",
}

// CreateOnDuplicateWithUpdateFields only update the no unique index field, it return error if updateFields contain unique index field.
// It doesn't update the Enabled and CreateTime field.
func (u *User) CreateOnDuplicateWithUpdateFields(db *gorm.DB, updateFields []string) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.CreateOnDuplicateWithUpdateFields",
	})()
	if len(updateFields) == 0 {
		return fmt.Errorf("Must have update fields.")
	}

	noUniqueIndexFields := (&UserDBFieldData{}).FetchNoneUniqueIndexFields()
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
			return fmt.Errorf("Field[%s] is unique index or wrong field or Enable field", UserDBFieldAndStructFieldRelate[field])
		}
		updateFieldsMap[field] = ""
	}

	if u.CreateTime.IsZero() {
		u.CreateTime = time.Now()
	}

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	u.Enabled = uint8(enumeration.BOOL__TRUE)

	structType := reflect.TypeOf(u).Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("Instance not struct type.")
	}
	structVal := reflect.ValueOf(u).Elem()

	var param_list []interface{}
	var str_list = []string{"insert into"}
	var insertFieldsStr = u.TableName() + "("
	var placeHolder = "values("
	for i := 0; i < structType.NumField(); i++ {
		if i == 0 {
			insertFieldsStr += UserStructFieldAndDBFieldRelate[structType.Field(i).Name]
			placeHolder += fmt.Sprintf("%s", "?")
		} else {
			insertFieldsStr += fmt.Sprintf(",%s", UserStructFieldAndDBFieldRelate[structType.Field(i).Name])
			placeHolder += fmt.Sprintf("%s", ", ?")
		}
		param_list = append(param_list, structVal.Field(i).Interface())
	}
	insertFieldsStr += ")"
	placeHolder += ")"
	str_list = append(str_list, []string{insertFieldsStr, placeHolder, "on duplicate key update"}...)

	var updateStr []string
	for i := 0; i < structType.NumField(); i++ {
		if dbField, ok := UserStructFieldAndDBFieldRelate[structType.Field(i).Name]; !ok {
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

func (u *User) DeleteById(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.DeleteById",
	})()

	err := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Delete(u).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordDeleteFailedError
	} else {
		return nil
	}
}

func (u *User) DeleteByPhone(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.DeleteByPhone",
	})()

	err := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Delete(u).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordDeleteFailedError
	} else {
		return nil
	}
}

func (u *User) DeleteByUserIDAndName(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.DeleteByUserIDAndName",
	})()

	err := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Delete(u).Error
	if err != nil {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordDeleteFailedError
	} else {
		return nil
	}
}

func (ul *UserList) FetchByCityID(db *gorm.DB, cityID uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByCityID",
	})()

	err := db.Table(User{}.TableName()).Where("F_city_id = ? and F_enabled = ?", cityID, enumeration.BOOL__TRUE).Find(ul).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (ul *UserList) FetchByCityIDAndAreaID(db *gorm.DB, cityID uint64, areaID int) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByCityIDAndAreaID",
	})()

	err := db.Table(User{}.TableName()).Where("F_city_id = ? and F_area_id = ? and F_enabled = ?", cityID, areaID, enumeration.BOOL__TRUE).Find(ul).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (u *User) FetchById(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchById",
	})()

	err := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Find(u).Error
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

func (u *User) FetchByIdForUpdate(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByIdForUpdate",
	})()

	err := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Set("gorm:query_option", "FOR UPDATE").Find(u).Error
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

func (u *User) FetchByPhone(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByPhone",
	})()

	err := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Find(u).Error
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

func (u *User) FetchByPhoneForUpdate(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByPhoneForUpdate",
	})()

	err := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Set("gorm:query_option", "FOR UPDATE").Find(u).Error
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

func (ul *UserList) FetchByUserID(db *gorm.DB, userID uint64) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByUserID",
	})()

	err := db.Table(User{}.TableName()).Where("F_user_id = ? and F_enabled = ?", userID, enumeration.BOOL__TRUE).Find(ul).Error
	if err == nil {
		return nil
	} else {
		logrus.Errorf("%s", err.Error())
		return dberr.RecordFetchFailedError
	}
}

func (u *User) FetchByUserIDAndName(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByUserIDAndName",
	})()

	err := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Find(u).Error
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

func (u *User) FetchByUserIDAndNameForUpdate(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchByUserIDAndNameForUpdate",
	})()

	err := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Set("gorm:query_option", "FOR UPDATE").Find(u).Error
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

func (ul *UserList) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.FetchList",
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

	err = db.Table(User{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("F_create_time desc").Find(ul).Error

	if err != nil {
		logrus.Errorf("%s", err.Error())
		return 0, dberr.RecordFetchFailedError
	} else {
		return int32(count), nil
	}
}

func (u *User) SoftDeleteById(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.SoftDeleteById",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = enumeration.BOOL__FALSE

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	err := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Updates(updateMap).Error
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
			delErr := db.Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Delete(&User{}).Error
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

func (u *User) SoftDeleteByPhone(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.SoftDeleteByPhone",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = enumeration.BOOL__FALSE

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	err := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Updates(updateMap).Error
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
			delErr := db.Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Delete(&User{}).Error
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

func (u *User) SoftDeleteByUserIDAndName(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.SoftDeleteByUserIDAndName",
	})()

	var updateMap = map[string]interface{}{}
	updateMap["F_enabled"] = enumeration.BOOL__FALSE

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	err := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Updates(updateMap).Error
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
			delErr := db.Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Delete(&User{}).Error
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

func (u *User) UpdateByIdWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.UpdateByIdWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = timelib.MySQLTimestamp(time.Now())

	}
	dbRet := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Updates(updateMap)
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
			findErr := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Find(&User{}).Error
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

func (u *User) UpdateByIdWithStruct(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.UpdateByIdWithStruct",
	})()

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	dbRet := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Updates(u)
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
			findErr := db.Table(u.TableName()).Where("F_id = ? and F_enabled = ?", u.Id, enumeration.BOOL__TRUE).Find(&User{}).Error
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

func (u *User) UpdateByPhoneWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.UpdateByPhoneWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = timelib.MySQLTimestamp(time.Now())

	}
	dbRet := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Updates(updateMap)
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
			findErr := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Find(&User{}).Error
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

func (u *User) UpdateByPhoneWithStruct(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.UpdateByPhoneWithStruct",
	})()

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	dbRet := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Updates(u)
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
			findErr := db.Table(u.TableName()).Where("F_phone = ? and F_enabled = ?", u.Phone, enumeration.BOOL__TRUE).Find(&User{}).Error
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

func (u *User) UpdateByUserIDAndNameWithMap(db *gorm.DB, updateMap map[string]interface{}) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.UpdateByUserIDAndNameWithMap",
	})()

	if _, ok := updateMap["F_update_time"]; !ok {
		updateMap["F_update_time"] = timelib.MySQLTimestamp(time.Now())

	}
	dbRet := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Updates(updateMap)
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
			findErr := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Find(&User{}).Error
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

func (u *User) UpdateByUserIDAndNameWithStruct(db *gorm.DB) error {
	defer duration.PrintDuration(map[string]interface{}{
		"request": "[DB]User.UpdateByUserIDAndNameWithStruct",
	})()

	if time.Time(u.UpdateTime).IsZero() {
		u.UpdateTime = timelib.MySQLTimestamp(time.Now())
	}

	dbRet := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Updates(u)
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
			findErr := db.Table(u.TableName()).Where("F_user_id = ? and F_name = ? and F_enabled = ?", u.UserID, u.Name, enumeration.BOOL__TRUE).Find(&User{}).Error
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
