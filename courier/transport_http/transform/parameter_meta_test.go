package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/courier/status_error"
	"github.com/johnnyeven/libtools/ptr"
	"github.com/johnnyeven/libtools/reflectx"
	"github.com/johnnyeven/libtools/timelib"
)

type JSON string

func (p JSON) MarshalJSON() ([]byte, error) {
	return []byte("json"), nil
}

func (p *JSON) UnmarshalJSON(data []byte) error {
	*p = JSON("json")
	return nil
}

type P string

func (p P) MarshalText() ([]byte, error) {
	return []byte("parameter"), nil
}

func (p *P) UnmarshalText(data []byte) error {
	*p = P(string(data))
	if string(*p) == "error" {
		return fmt.Errorf("error")
	}
	*p = P("parameter")
	return nil
}

type Uint64List []uint64

func (list Uint64List) MarshalJSON() ([]byte, error) {
	if len(list) == 0 {
		return []byte(`[]`), nil
	}
	strValues := make([]string, 0)
	for _, v := range list {
		strValues = append(strValues, fmt.Sprintf(`"%d"`, v))
	}
	return []byte(`[` + strings.Join(strValues, ",") + `]`), nil
}

func (list *Uint64List) UnmarshalJSON(data []byte) (err error) {
	strValues := make([]string, 0)
	err = json.Unmarshal(data, &strValues)
	if err != nil {
		return err
	}
	finalList := Uint64List{}
	for i, strValue := range strValues {
		v, parseErr := strconv.ParseUint(strValue, 10, 64)
		if parseErr != nil {
			err = fmt.Errorf(`[%d] cannot unmarshal string into value of type uint64`, i)
			return
		}
		finalList = append(finalList, v)
	}
	*list = finalList
	return
}

func TestParameterMeta(t *testing.T) {
	tt := assert.New(t)

	type Data struct {
		String  string  `json:"string" validate:"@string[3,]"`
		Pointer *string `json:"pointer" validate:"@string[3,]"`
	}

	type SomeReq struct {
		P                   P                      `name:"p" in:"query"`
		Slice               []string               `name:"slice" in:"query"`
		Array               [5]string              `name:"array" in:"query"`
		Query               string                 `name:"query" in:"query" validate:"@string[3,]"`
		Uint64List          Uint64List             `name:"Uint64List" in:"query" default:"" validate:"@array[1,]:@uint64[3,]"`
		QueryWithDefaults   string                 `name:"queryWithDefaults" in:"query" default:"123123" validate:"@string[3,]"`
		Pointer             *string                `name:"pointerWithDefaults" in:"path" default:"123123" validate:"@string[3,]"`
		PointerWithDefaults *string                `name:"pointer" in:"path" validate:"@string[3,]"`
		Bytes               []byte                 `name:"bytes" in:"query" style:"form,explode" `
		CreateTime          timelib.MySQLTimestamp `name:"createTime" in:"query" default:""`
		UpdateTime          timelib.MySQLTimestamp `name:"updateTime,omitempty" in:"query"`
		Data                Data                   `in:"body"`
		PtrData             *Data                  `in:"body"`
		DataSlice           []string               `in:"body"`
		DataArray           [5]string              `in:"body"`
	}

	req := &SomeReq{}
	tpe := reflectx.IndirectType(reflect.TypeOf(req))
	rv := reflect.Indirect(reflect.ValueOf(req))

	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		fieldValue := rv.Field(i)

		tagIn, _, tagInFlags := GetTagIn(&field)

		p := NewParameterMeta(&field, fieldValue, tagIn, tagInFlags)

		switch p.Name {
		case "Uint64List":
			{
				err := p.UnmarshalStringAndValidate()
				tt.Equal(status_error.ErrorFields{
					{
						In:    "query",
						Field: "Uint64List",
						Msg:   "切片元素个数不在[1， 0]范围内，当前个数：0",
					},
				}, status_error.FromError(err).ErrorFields)
			}
			{
				err := p.UnmarshalStringAndValidate("123", "123")
				tt.NoError(err)
				tt.Equal(Uint64List{123, 123}, p.Value.Interface())
			}
		case "queryWithDefaults":
			{
				err := p.UnmarshalStringAndValidate("")
				tt.NoError(err)
				tt.Equal("123123", p.Value.Interface())
			}
		case "pointerWithDefaults":
			{
				err := p.UnmarshalStringAndValidate("")
				tt.NoError(err)
				tt.Equal(ptr.String("123123"), p.Value.Interface())
			}
		case "query":
			{
				err := p.UnmarshalStringAndValidate("")
				tt.NotNil(err)
				tt.Equal(status_error.ErrorFields{
					{
						In:    "query",
						Field: "query",
						Msg:   ErrMsgForRequired,
					},
				}, status_error.FromError(err).ErrorFields)
			}

			{
				err := p.UnmarshalStringAndValidate("100")
				tt.Nil(err)
				tt.Equal("100", p.Value.Interface())
			}
		case "p":
			{
				err := p.UnmarshalStringAndValidate("error")
				tt.Error(err)
				tt.Equal(status_error.ErrorFields{
					{
						In:    "query",
						Field: "p",
						Msg:   "error",
					},
				}, status_error.FromError(err).ErrorFields)
			}
			{
				err := p.UnmarshalStringAndValidate("100")
				tt.Nil(err)
				tt.Equal(P("parameter"), p.Value.Interface())
			}
		case "pointer":
			{
				err := p.UnmarshalStringAndValidate("")
				tt.Nil(p.Value.Interface())
				tt.NotNil(err)
				tt.Equal(status_error.ErrorFields{
					status_error.NewErrorField("path", "pointer", ErrMsgForRequired),
				}, status_error.FromError(err).ErrorFields.Sort())
			}
			{
				err := p.UnmarshalStringAndValidate("10")
				tt.Equal(ptr.String("10"), p.Value.Interface())
				tt.Error(err)
				tt.Equal(status_error.ErrorFields{
					status_error.NewErrorField("path", "pointer", "字符串长度不在[3， 1024]范围内，当前长度：2"),
				}, status_error.FromError(err).ErrorFields.Sort())
			}
			{
				err := p.UnmarshalStringAndValidate("100")
				tt.Equal(ptr.String("100"), p.Value.Interface())
				tt.NoError(err)
			}
		case "bytes":
			{
				err := p.UnmarshalStringAndValidate("111")
				tt.NoError(err)
				tt.Equal([]byte("111"), p.Value.Interface())
			}
		case "slice":
			{
				err := p.UnmarshalStringAndValidate("111", "222")
				tt.NoError(err)
				tt.Equal([]string{"111", "222"}, p.Value.Interface())
			}
		case "array":
			{
				err := p.UnmarshalStringAndValidate("111", "222")
				tt.NoError(err)
				tt.Equal([5]string{"111", "222", "", "", ""}, p.Value.Interface())
			}
		case "updateTime":
			{
				err := p.UnmarshalStringAndValidate()
				tt.NoError(err)
			}
		case "createTime":
			{
				err := p.UnmarshalStringAndValidate("2017-10-10T00:00:00Z")
				tt.Nil(err)
				d, _ := timelib.ParseMySQLTimestampFromString("2017-10-10T00:00:00Z")
				tt.Equal(d, p.Value.Interface())
			}
		case "PtrData":
			{
				buf := bytes.NewBufferString(`{"string":"1"}`)
				err := p.UnmarshalFromReader(buf)
				tt.NotNil(err)
				tt.Equal(status_error.ErrorFields{
					status_error.NewErrorField("body", "pointer", ErrMsgForRequired),
					status_error.NewErrorField("body", "string", "字符串长度不在[3， 1024]范围内，当前长度：1"),
				}, status_error.FromError(err).ErrorFields.Sort())
			}
		case "Data":
			{
				err := p.UnmarshalFromReader(nil)
				tt.NotNil(err)
				tt.Equal(int64(status_error.ReadFailed), status_error.FromError(err).Code)
			}

			{
				file, _ := ioutil.TempFile("", "")
				file.Close()
				err := p.UnmarshalFromReader(file)
				tt.NotNil(err)
				tt.Equal(int64(status_error.ReadFailed), status_error.FromError(err).Code)
			}

			{
				file, _ := ioutil.TempFile("", "")
				err := p.UnmarshalFromReader(file)
				tt.NotNil(err)
				tt.Equal(int64(status_error.InvalidBodyStruct), status_error.FromError(err).Code)
				file.Close()
			}
			{
				buf := bytes.NewBufferString(`{"string":"1"}`)
				err := p.UnmarshalFromReader(buf)
				tt.NotNil(err)
				tt.Equal(status_error.ErrorFields{
					status_error.NewErrorField("body", "pointer", ErrMsgForRequired),
					status_error.NewErrorField("body", "string", "字符串长度不在[3， 1024]范围内，当前长度：1"),
				}, status_error.FromError(err).ErrorFields.Sort())
			}
			{
				buf := bytes.NewBufferString(`{"string":"111", "pointer":1}`)
				err := p.UnmarshalFromReader(buf)
				tt.NotNil(err)
				tt.Equal(status_error.ErrorFields{
					status_error.NewErrorField("body", "pointer", "json: cannot unmarshal number into Go struct field Data.pointer of type string"),
				}, status_error.FromError(err).ErrorFields.Sort())
			}
			{
				buf := bytes.NewBufferString(`{"string":"111","pointer":"111"}`)
				err := p.UnmarshalFromReader(buf)
				tt.Nil(err)
				tt.Equal(Data{
					String:  "111",
					Pointer: ptr.String("111"),
				}, p.Value.Interface())
			}
		case "DataSlice":
			{
				buf := bytes.NewBufferString(`["123","123"]`)
				err := p.UnmarshalFromReader(buf)
				tt.Nil(err)
				tt.Equal([]string{
					"123",
					"123",
				}, p.Value.Interface())
			}
		case "DataArray":
			{
				buf := bytes.NewBufferString(`["123","123"]`)
				err := p.UnmarshalFromReader(buf)
				tt.Nil(err)
				tt.Equal([5]string{
					"123",
					"123",
					"",
					"",
					"",
				}, p.Value.Interface())
			}
		}
	}
}

func TestParameterMeta_Marshal(t *testing.T) {
	tt := assert.New(t)

	type Data struct {
		String string `json:"string" validate:"@string[3,]"`
	}

	type SomeReq struct {
		JSON       JSON                   `name:"json" in:"query"`
		P          P                      `name:"p" in:"query"`
		Bytes      []byte                 `name:"bytes" in:"query" style:"form,explode" `
		Slice      []string               `name:"slice" in:"query"`
		Query      string                 `name:"query" in:"query" validate:"@string[3,]"`
		Pointer    *string                `name:"pointer" in:"path" validate:"@string[3,]"`
		Pointer2   *string                `name:"pointerIgnore" in:"path" validate:"@string[3,]"`
		CreateTime timelib.MySQLTimestamp `name:"createTime" in:"query" default:""`
		Data       Data                   `in:"body"`
		DataSlice  []string               `in:"body"`
		FormData   string                 `in:"formData"`
	}

	req := &SomeReq{
		P:       "!",
		Query:   "query",
		Pointer: ptr.String("pointer"),
		Bytes:   []byte("bytes"),
		Slice: []string{
			"1", "2",
		},
		Data: Data{
			String: "string",
		},
		DataSlice: []string{
			"1", "2",
		},
		FormData:   "1",
		CreateTime: timelib.MySQLTimestamp(timelib.Now()),
	}

	tpe := reflectx.IndirectType(reflect.TypeOf(req))
	rv := reflect.Indirect(reflect.ValueOf(req))

	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		fieldValue := rv.Field(i)

		tagIn, _, tagInFlags := GetTagIn(&field)

		p := NewParameterMeta(&field, fieldValue, tagIn, tagInFlags)

		switch p.Name {
		case "json":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("json"), dataList)
		case "query":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("query"), dataList)
		case "p":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("parameter"), dataList)
		case "pointer":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("pointer"), dataList)
		case "createTime":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList(req.CreateTime.String()), dataList)
		case "slice":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("1", "2"), dataList)
		case "pointerIgnore":
			dataList, _ := p.Marshal()
			tt.Nil(dataList)
		case "bytes":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("bytes"), dataList)
		case "Data":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList(`{"string":"string"}`), dataList)
		case "DataSlice":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList(`["1","2"]`), dataList)
		case "FormData":
			dataList, _ := p.Marshal()
			tt.Equal(BytesList("1"), dataList)
		}
	}
}
