package gen_method

import (
	"fmt"
	"go/build"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/johnnyeven/libtools/codegen"
	"github.com/johnnyeven/libtools/codegen/loaderx"
)

type Field struct {
	Name                string
	Type                string
	DbFieldName         string
	IndexName           string
	IndexNumber         int
	IsEnable            bool
	IsCreateTime        bool
	IsSpecifyIndexOrder bool
}

type Model struct {
	Pkg                 *types.Package
	Name                string
	Fields              []*Field
	IsDbModel           bool
	UniqueIndex         map[string][]Field
	NormalIndex         map[string][]Field
	PrimaryIndex        []Field
	HasCreateTimeField  bool
	HasUpdateTimeField  bool
	HasEnabledField     bool
	EnabledFieldType    string
	CreateTimeFieldType string
	UpdateTimeFieldType string
	DbCreateTimeField   string
	DbEnabledField      string
	DbUpdateTimeField   string
	FuncMapContent      map[string]string
	Deps                []string
}

func (model *Model) collectInfoFromStructType(tpeStruct *types.Struct) {
	for i := 0; i < tpeStruct.NumFields(); i++ {
		field := tpeStruct.Field(i)
		tag := reflect.StructTag(tpeStruct.Tag(i))

		sqlSettings := ParseTagSetting(tag.Get("sql"))
		gormSettings := ParseTagSetting(tag.Get("gorm"))

		if len(gormSettings) != 0 && len(sqlSettings) != 0 {
			model.IsDbModel = true
		} else {
			continue
		}

		var tmpField = Field{}
		if dbFieldName, ok := gormSettings["COLUMN"]; ok {
			tmpField.DbFieldName = dbFieldName[0]
		} else {
			tmpField.DbFieldName = tmpField.Name
		}
		tmpField.Name = field.Name()

		pkgPath, method := loaderx.GetPkgImportPathAndExpose(field.Type().String())
		if pkgPath != "" {
			pkg, _ := build.Import(pkgPath, "", build.ImportComment)
			tmpField.Type = fmt.Sprintf("%s.%s", pkg.Name, method)
		} else {
			tmpField.Type = method
		}

		if pkgPath != "" {
			model.Deps = append(model.Deps, pkgPath)
		}

		model.Fields = append(model.Fields, &tmpField)
		if tmpField.Name == "Enabled" || tmpField.DbFieldName == "F_enabled" {
			// enabled field don't join into index slice
			tmpField.IsEnable = true
			model.HasEnabledField = true
			model.EnabledFieldType = tmpField.Type
			model.DbEnabledField = tmpField.DbFieldName
			continue
		} else if tmpField.Name == "UpdateTime" || tmpField.DbFieldName == "F_update_time" {
			model.HasUpdateTimeField = true
			model.UpdateTimeFieldType = tmpField.Type
			model.DbUpdateTimeField = tmpField.DbFieldName
		} else if tmpField.Name == "CreateTime" || tmpField.DbFieldName == "F_create_time" {
			tmpField.IsCreateTime = true
			model.HasCreateTimeField = true
			model.CreateTimeFieldType = tmpField.Type
			model.DbCreateTimeField = tmpField.DbFieldName
		}

		if _, ok := gormSettings["PRIMARY_KEY"]; ok {
			model.PrimaryIndex = append(model.PrimaryIndex, tmpField)
		}
		if indexName, ok := sqlSettings["INDEX"]; ok && len(indexName) > 0 {
			for _, index := range indexName {
				tmpField.IndexName, tmpField.IndexNumber, tmpField.IsSpecifyIndexOrder = ParseIndex(index)
				if _, ok := model.NormalIndex[tmpField.IndexName]; ok {
					model.NormalIndex[tmpField.IndexName] = append(
						model.NormalIndex[tmpField.IndexName],
						tmpField)
				} else {
					model.NormalIndex[tmpField.IndexName] = []Field{tmpField}
				}
			}
		}
		if indexName, ok := sqlSettings["UNIQUE_INDEX"]; ok && len(indexName) > 0 {
			for _, index := range indexName {
				tmpField.IndexName, tmpField.IndexNumber, tmpField.IsSpecifyIndexOrder = ParseIndex(index)
				if _, ok := model.UniqueIndex[tmpField.IndexName]; ok {
					model.UniqueIndex[tmpField.IndexName] = append(
						model.UniqueIndex[tmpField.IndexName],
						tmpField)
				} else {
					model.UniqueIndex[tmpField.IndexName] = []Field{tmpField}
				}
			}
		}
	}
}

func (model *Model) Output(pkgName string, ignoreCreateTableNameFunc bool) {
	if err := genTableNameFunc(model, pkgName, "golib", ignoreCreateTableNameFunc); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genCreateFunc(model); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	model.genDBFieldFunc()
	model.genCodeByNormalIndex()
	model.genCodeByUniqueIndex()
	model.genCodeByPrimaryKeyIndex()

	if err := genFetchListFunc(model); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	model.GenerateFile()
}

func (model *Model) GenerateFile() {
	var funcNameList []string

	// first part of file
	var tableName = "tableName"
	for key := range model.FuncMapContent {
		if key == tableName {
			continue
		}
		funcNameList = append(funcNameList, key)
	}

	contents := []string{
		model.FuncMapContent[tableName],
	}

	sort.Strings(funcNameList)

	for _, funcName := range funcNameList {
		contents = append(contents, model.FuncMapContent[funcName])
	}

	p, _ := build.Import(model.Pkg.Path(), "", build.FindOnly)
	cwd, _ := os.Getwd()
	path, _ := filepath.Rel(cwd, p.Dir)

	filename := path + "/" + replaceUpperWithLowerAndUnderscore(model.Name) + ".go"
	content := strings.Join(contents, "\n\n")
	bytes, err := imports.Process(filename, []byte(content), nil)
	if err != nil {
		panic(err)
	} else {
		content = string(bytes)
	}
	codegen.WriteFile(codegen.GeneratedSuffix(filename), content)
}

func (model *Model) genCodeByUniqueIndex() {
	for _, fieldList := range model.UniqueIndex {
		sortFieldList := model.sortFieldsByIndexNumber(fieldList)
		model.handleGenCodeForUniqueIndex(sortFieldList)
	}
}

func (model *Model) genCodeByPrimaryKeyIndex() {
	if len(model.PrimaryIndex) > 0 {
		model.handleGenCodeForUniqueIndex(model.PrimaryIndex)
	}
}

func (model *Model) genCodeByNormalIndex() {
	for _, fieldList := range model.NormalIndex {
		sortFieldList := model.sortFieldsByIndexNumber(fieldList)
		baseInfoGenCode := fetchBaseInfoOfGenFuncForNormalIndex(sortFieldList)
		if err := genFetchFuncByNormalIndex(model, baseInfoGenCode); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		model.handleGenFetchCodeBySubIndex(sortFieldList)
	}
}

func (model *Model) genDBFieldFunc() {
	var structFields, structFieldsAndValue, noUniqueIndexFields, dBAndStructField, structAndDbField []string
	var primaryIndexMap = make(map[string]string)
	for _, field := range model.PrimaryIndex {
		primaryIndexMap[field.Name] = ""
	}

	for _, field := range model.Fields {
		structFields = append(structFields, fmt.Sprintf("%s string", field.Name))
		structFieldsAndValue = append(structFieldsAndValue, fmt.Sprintf("%s:\"%s\"", field.Name, field.DbFieldName))
		dBAndStructField = append(dBAndStructField, fmt.Sprintf("\"%s\" : \"%s\",", field.DbFieldName, field.Name))
		structAndDbField = append(structAndDbField, fmt.Sprintf("\"%s\" : \"%s\",", field.Name, field.DbFieldName))
		if _, ok := model.UniqueIndex[field.IndexName]; !ok {
			if _, ok := primaryIndexMap[field.Name]; !ok && !field.IsEnable && !field.IsCreateTime {
				noUniqueIndexFields = append(noUniqueIndexFields, fmt.Sprintf("\"%s\"", field.DbFieldName))
			}
		}
	}

	if err := genDBFiledFunc(model, structFields, structFieldsAndValue, noUniqueIndexFields, dBAndStructField,
		structAndDbField); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func (model *Model) sortFieldsByIndexNumber(fields []Field) []Field {
	if len(fields) == 0 && len(fields) == 1 {
		return fields
	}

	if !IsSpecifyIndexSequence(fields) {
		return fields
	}

	var markIndexNubmer = make(map[string]string)
	var sortFieldSlice = []Field{}
	for _, field := range fields {
		if fieldName, ok := markIndexNubmer[fmt.Sprintf("%d", field.IndexNumber)]; ok {
			fmt.Printf("Field[%s] and Field[%s] has same index number[%d] in Mode[%s]", fieldName, field.Name,
				field.IndexNumber, model.Name)
			os.Exit(1)
		}
		if field.IndexNumber >= 0 {
			tmpSlice := make([]Field, field.IndexNumber+1)
			tmpSlice[field.IndexNumber] = field
			if len(sortFieldSlice) > field.IndexNumber+1 {
				sortFieldSlice = append(tmpSlice, sortFieldSlice[field.IndexNumber+1:]...)
			} else if len(sortFieldSlice) < field.IndexNumber+1 {
				sortFieldSlice = append(sortFieldSlice, tmpSlice[len(sortFieldSlice):]...)
			} else {
				fmt.Printf("Field[%s] wrong index sequence, may be same index number.\n", field.Name)
				os.Exit(1)
			}
		}
	}

	var notEmptySlice = []Field{}
	for index, field := range sortFieldSlice {
		if len(field.Name) > 0 {
			notEmptySlice = append(notEmptySlice, sortFieldSlice[index])
			//sortFieldSlice = append(sortFieldSlice[:index], sortFieldSlice[index+1:]...)
		}
	}

	return notEmptySlice
}

func (model *Model) genBatchFetchFuncBySingleIndex(field Field) {
	if err := genBatchFetchFunc(model, field.Name, field.DbFieldName, field.Type); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func (model *Model) handleGenFetchCodeBySubIndex(fieldList []Field) {
	if len(fieldList) == 1 {
		model.genBatchFetchFuncBySingleIndex(fieldList[0])
	} else if len(fieldList) > 1 {
		// [x, y, z, e] Split to [x, y, z], [x, y], [x]
		for i := 1; i < len(fieldList); i++ {
			subSortFieldSlice := fieldList[:len(fieldList)-i]
			baseInfoGenCode := fetchBaseInfoOfGenFuncForNormalIndex(subSortFieldSlice)
			if err := genFetchFuncByNormalIndex(model, baseInfoGenCode); err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}

			if len(subSortFieldSlice) == 1 {
				model.genBatchFetchFuncBySingleIndex(subSortFieldSlice[0])
			}

		}
	}
}

func (model *Model) handleGenCodeForUniqueIndex(sortFieldList []Field) {
	baseInfoGenCode := fetchBaseInfoOfGenFuncForUniqueIndex(model, sortFieldList)
	if err := genFetchFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	model.handleGenFetchCodeBySubIndex(sortFieldList)

	if err := genFetchForUpdateFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genUpdateWithStructFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genUpdateWithMapFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genSoftDeleteFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	if err := genPhysicsDeleteFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
