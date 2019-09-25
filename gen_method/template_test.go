package gen_method

import (
	"os"
	"testing"
	"text/template"

	"github.com/johnnyeven/libtools/codegen"
)

var gFile = os.Stdout

func TestTableNameTemplate(t *testing.T) {
	t.Log(gFile)
	data := TableNameTemplateParam{
		StructName:   "UserTest",
		ReceiverName: "ut",
		PackageName:  "modules",
	}

	tmpl := template.Must(template.New("create").Parse(TableNameTemplate))
	if err := tmpl.Execute(gFile, data); err != nil {
		t.Error(err)
	}
}

func TestCreateTemplate(t *testing.T) {
	t.Log(gFile)
	data := TableNameTemplateParam{
		StructName:   "UserTest",
		ReceiverName: "ut",
	}

	tmpl := template.Must(template.New("create").Parse(CreateTemplate))
	if err := tmpl.Execute(gFile, data); err != nil {
		t.Error(err)
	}
}

func TestFetchTemplate(t *testing.T) {
	data := new(FetchTemplateParam)
	data.StructName = "UserTest"
	data.ReceiverName = "ut"
	data.Field = "Id"
	data.DbField = "F_id"

	tmpl := template.Must(template.New("create").Parse(FetchTemplate))
	if err := tmpl.Execute(gFile, data); err != nil {
		t.Error(err)
	}
}

func TestBatchFetchTemplate(t *testing.T) {
	data := new(BatchFetchTemplateParam)
	data.StructName = "UserTest"
	data.ReceiverName = "ut"
	data.Field = "Id"
	data.DbField = "F_id"
	data.FieldType = "uint64"

	tmpl := template.Must(template.New("create").Parse(BatchFetchTemplate))
	if err := tmpl.Execute(gFile, data); err != nil {
		t.Error(err)
	}
}

func TestFetchList(t *testing.T) {
	data := FetchListTemplateParam{}
	data.StructName = "UserTest"
	data.ReceiverListName = "utl"
	data.StructListName = "UserTestList"
	tmpl := template.Must(template.New("create").Parse(FetchListTemplate))
	if err := tmpl.Execute(gFile, data); err != nil {
		t.Error(err)
	}

}

func TestSplitStringByWord(t *testing.T) {
	tmpStr := []string{
		"adcIDTest",
		"AbcDref",
		"bbIDArtisan",
		"GenerateURITest",
		"GolangURIIDIPTest",
		"test123",
		"niceGirl",
		"GolangURIIDIPTestAPIIDasebb",
	}

	t.Log(tmpStr[0][len(tmpStr[0]):])
	for _, str := range tmpStr {
		t.Log(str, codegen.ToUpperSnakeCase(str))
	}
}

func TestMain(m *testing.M) {
	//file, err := os.Create("./generate_user.go")
	//if err != nil {
	//panic(err)
	//}
	//gFile = file
	os.Exit(m.Run())
}
