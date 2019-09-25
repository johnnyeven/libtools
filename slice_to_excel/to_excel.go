package slice_to_excel

import (
	"fmt"
	"reflect"

	"github.com/tealeg/xlsx"

	"github.com/johnnyeven/libtools/timelib"
)

type XLSXLabel interface {
	Label() string
}

func GetExcel(sheetName string, srcSlice interface{}) (file *xlsx.File, err error) {
	if reflect.TypeOf(srcSlice).Kind() != reflect.Slice {
		panic(fmt.Errorf("srcSlice must be an slice"))
	}

	file = xlsx.NewFile()

	err = SetExcel(file, sheetName, reflect.ValueOf(srcSlice))

	return
}

func SetExcel(f *xlsx.File, sheetName string, reflectValue reflect.Value) (err error) {
	sheet, _ := f.AddSheet(sheetName)
	row := sheet.AddRow()
	reflectType := reflectValue.Type().Elem()

	setLabelRow(row, reflectType)

	for i := 0; i < reflectValue.Len(); i++ {
		row := sheet.AddRow()
		err := setToExcelRow(row, reflectValue.Index(i))
		if err != nil {
			return err
		}
	}

	return err
}

func setLabelRow(row *xlsx.Row, reflectType reflect.Type) {
	if reflectType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("item is not struct"))
	}

	for i := 0; i < reflectType.NumField(); i++ {
		structField := reflectType.Field(i)

		if structField.Anonymous {
			setLabelRow(row, structField.Type)
			continue
		}
		xlsTag := structField.Tag.Get("xlsx")
		if xlsTag != "" {
			cell := row.AddCell()
			cell.SetValue(xlsTag)
		}
	}
}

func setToExcelRow(row *xlsx.Row, reflectValue reflect.Value) (err error) {
	reflectValue = reflect.Indirect(reflectValue)
	reflectType := reflectValue.Type()

	for i := 0; i < reflectType.NumField(); i++ {
		structField := reflectType.Field(i)
		structFieldValue := reflectValue.Field(i)

		if structField.Anonymous {
			err = setToExcelRow(row, structFieldValue)
			continue
		}

		xlsTag := structField.Tag.Get("xlsx")
		if xlsTag != "" {
			cell := row.AddCell()
			formatTag := structField.Tag.Get("format")
			if formatTag != "" {
				if t, ok := structFieldValue.Interface().(int64); ok {
					if formatTag == "currency" {
						cell.SetFloat(float64(t) / 100)
					}
				}
			} else {
				if t, ok := structFieldValue.Interface().(timelib.MySQLTimestamp); ok {
					if !t.IsZero() {
						cell.SetString(t.String())
					} else {
						cell.SetString("")
					}
				} else if t, ok := structFieldValue.Interface().(XLSXLabel); ok {
					cell.SetString(t.Label())
				} else if t, ok := structFieldValue.Interface().(fmt.Stringer); ok {
					cell.SetString(t.String())
				} else {
					switch structFieldValue.Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
						cell.SetValue(structFieldValue.Interface())
					case reflect.Int64, reflect.Uint64:
						cell.SetString(fmt.Sprintf("%d", structFieldValue.Interface()))
					case reflect.Float32:
						cell.SetFloat(float64(structFieldValue.Interface().(float32)))
					case reflect.Float64:
						cell.SetFloat(structFieldValue.Interface().(float64))
					case reflect.Bool:
						cell.SetBool(structFieldValue.Interface().(bool))
					default:
						panic(fmt.Sprintf("field[%s] type [%s] is not support", structField.Name, structFieldValue.Kind()))
					}
				}
			}
		}
	}
	return
}
