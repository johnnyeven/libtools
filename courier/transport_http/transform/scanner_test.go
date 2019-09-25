package transform

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/johnnyeven/libtools/timelib"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/ptr"
)

type Gender int

const (
	GenderMale Gender = iota + 1
	GenderFemale
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	}
	return ""
}

func (g Gender) MarshalJSON() ([]byte, error) {
	return []byte(g.String()), nil
}

func (g *Gender) UnmarshalJSON(data []byte) error {
	s, _ := strconv.Unquote(string(data))
	switch s {
	case "male":
		*g = GenderMale
	case "female":
		*g = GenderMale
	default:
		return fmt.Errorf("unknown gender")
	}
	return nil
}

func TestScanner(t *testing.T) {
	tt := assert.New(t)

	type Sub struct {
		Name string `json:"name,omitempty"`
	}

	type TestModel struct {
		Gender              Gender   `json:"gender" default:""`
		Slice               []string `json:"slice"`
		PtrString           *string  `json:"ptrString" default:""`
		PtrStringUseInput   *string  `json:"ptrStringUseInput" default:"2"`
		PtrStringUseDefault *string  `json:"ptrStringUseDefault" default:"2"`
		RequiredString      *string  `json:"requiredString"`
		Struct              Sub      `json:"struct"`
		PtrStruct           *Sub     `json:"ptrStruct,omitempty"`
	}

	cases := []struct {
		desc       string
		model      TestModel
		finalModel TestModel
		valid      bool
		msgMap     ErrMsgMap
	}{
		{
			desc:  "normal",
			valid: true,
			model: TestModel{
				PtrStringUseInput: ptr.String(""),
				RequiredString:    ptr.String(""),
			},
			finalModel: TestModel{
				PtrStringUseInput:   ptr.String(""),
				PtrStringUseDefault: ptr.String("2"),
				RequiredString:      ptr.String(""),
			},
		},
	}

	for _, tc := range cases {
		tpe := reflect.TypeOf(tc.model)
		rv := reflect.Indirect(reflect.ValueOf(&tc.model))

		valid, msgMap := NewScanner().Validate(rv, tpe)

		tt.Equal(tc.valid, valid, tc.desc)
		tt.Equal(tc.msgMap, msgMap, tc.desc)
		tt.Equal(tc.finalModel, tc.model, tc.desc)
	}
}

type TestSubModel1 struct {
	String string          `json:"string" default:"" validate:"@string[1,]"`
	Slice  []TestSubModel2 `json:"slice"`
}

type TestSubModel2 struct {
	String string `json:"string" default:"" validate:"@string[1,]"`
}

type TestModel2 struct {
	a        string
	Gender   Gender                 `json:"gender" default:"trans"`
	SliceLen []string               `json:"slice_len" validate:"@array[1,]"`
	Slice    []string               `json:"slice" validate:"@array[1,]:@string[1,]"`
	Hook     string                 `json:"hook"`
	Object   TestSubModel1          `json:"object"`
	Time     timelib.MySQLTimestamp `json:"time"`
}

func (v TestModel2) ValidateHook() string {
	return "hook failed"
}

func TestScannerWithValidateErrors(t *testing.T) {
	tt := assert.New(t)

	cases := []struct {
		desc       string
		model      TestModel2
		finalModel TestModel2
		valid      bool
		msgMap     ErrMsgMap
	}{
		{
			desc:  "default set failed",
			valid: false,
			model: TestModel2{
				Slice: []string{"", ""},
				Object: TestSubModel1{
					Slice: []TestSubModel2{
						{
							String: "",
						},
					},
				},
			},
			finalModel: TestModel2{
				Slice: []string{"", ""},
				Object: TestSubModel1{
					Slice: []TestSubModel2{
						{
							String: "",
						},
					},
				},
			},
			msgMap: ErrMsgMap{
				"gender":                 "Gender can't set wrong default value trans",
				"slice_len":              "切片元素个数不在[1， 0]范围内，当前个数：0",
				"slice":                  "切片元素不满足校验[字符串长度不在[1， 1024]范围内，当前长度：0]",
				"hook":                   "hook failed",
				"time":                   "缺失必填字段",
				"object.string":          "字符串长度不在[1， 1024]范围内，当前长度：0",
				"object.slice[0].string": "字符串长度不在[1， 1024]范围内，当前长度：0",
			},
		},
	}

	for _, tc := range cases {
		tpe := reflect.TypeOf(tc.model)
		rv := reflect.Indirect(reflect.ValueOf(&tc.model))

		valid, msgMap := NewScanner().Validate(rv, tpe)

		tt.Equal(tc.valid, valid, tc.desc)
		tt.Equal(tc.msgMap, msgMap, tc.desc)
		tt.Equal(tc.finalModel, tc.model, tc.desc)
	}
}

func TestMarshalAndValidate(t *testing.T) {
	tt := assert.New(t)

	{
		cases := []struct {
			desc             string
			v                int
			defaultValue     string
			required         bool
			tagValidate      string
			tagErrMsg        string
			errMsg           string
			needToSetDefault bool
		}{
			{
				desc:     "int required",
				required: true,
				errMsg:   ErrMsgForRequired,
			},
			{
				desc:             "set int with default value",
				defaultValue:     "1",
				needToSetDefault: true,
			},
			{
				desc:         "set int with default value failed",
				defaultValue: "str",
				errMsg:       fmt.Sprintf("%s can't set wrong default value %s", "int", "str"),
			},
		}

		for _, tc := range cases {
			tpe := reflect.TypeOf(tc.v)
			rv := reflect.Indirect(reflect.ValueOf(&tc.v))

			tt.Equal(tc.errMsg, MarshalAndValidate(rv, tpe, tc.defaultValue, tc.required, tc.tagValidate, tc.tagErrMsg), tc.desc)
			if tc.needToSetDefault {
				tt.Equal(tc.defaultValue, fmt.Sprintf("%v", tc.v))
			}
		}
	}

	{
		cases := []struct {
			desc             string
			V                *string
			defaultValue     string
			required         bool
			tagValidate      string
			tagErrMsg        string
			errMsg           string
			needToSetDefault bool
		}{
			{
				desc: "nil string skip",
			},
			{
				desc:     "nil string required",
				required: true,
				errMsg:   ErrMsgForRequired,
			},
			{
				desc:             "nil string should set default",
				defaultValue:     "1",
				needToSetDefault: true,
			},
			{
				desc:             "nil string should set default",
				defaultValue:     "1",
				needToSetDefault: true,
			},
			{
				desc:             "ptr string should key value",
				V:                ptr.String(""),
				defaultValue:     "1",
				needToSetDefault: false,
			},
		}

		for _, tc := range cases {
			tpe := reflect.TypeOf(tc.V)
			rv := reflect.Indirect(reflect.ValueOf(&tc)).FieldByName("V")

			tt.Equal(tc.errMsg, MarshalAndValidate(rv, tpe, tc.defaultValue, tc.required, tc.tagValidate, tc.tagErrMsg), tc.desc)

			if tc.needToSetDefault {
				tt.Equal(&tc.defaultValue, tc.V, tc.desc)
			}
		}
	}

}
