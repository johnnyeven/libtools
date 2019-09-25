package conf_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/johnnyeven/libtools/conf/presets"
	"github.com/johnnyeven/libtools/courier/transport_http/transform"
	"github.com/johnnyeven/libtools/ptr"

	"github.com/johnnyeven/libtools/conf"

	"github.com/stretchr/testify/assert"
)

type E struct {
	E int
}

func (e E) MarshalDefaults(v interface{}) {
	if target, ok := v.(*E); ok {
		if target.E == 0 {
			target.E = 4
		}
	}
}

type D struct {
	E
	F float64
	G string
}

type StructStruct struct {
	IntPointer *int
	A          int
	B          *bool
	C          uint
	D          D
	DP         *D
	List       []string
	Timeout    time.Duration
}

func TestUnmarshal(t *testing.T) {
	config := &StructStruct{
		A: 1,
		B: ptr.Bool(false),
		C: 2,
		D: D{
			E: E{},
		},
		DP: &D{
			E: E{
				E: 4,
			},
			F: 2.0,
			G: "abc",
		},
		List: []string{"def"},
	}

	os.Setenv("TEST_INTPOINTER", "0")
	os.Setenv("TEST_A", "2")
	os.Setenv("TEST_B", "true")
	os.Setenv("TEST_D_E", "5")
	os.Setenv("TEST_DP_E", "5")
	os.Setenv("TEST_LIST", "gds,123")
	os.Setenv("TEST_TIMEOUT", "5s")

	ok, errMsgs := conf.NewScanner("TEST").Unmarshal(reflect.ValueOf(config), reflect.TypeOf(config))
	if !ok {
		t.Log(errMsgs)
	}

	assert.Equal(t, &StructStruct{
		A: 2,
		B: ptr.Bool(true),
		C: 2,
		D: D{
			E: E{
				E: 5,
			},
		},
		DP: &D{
			E: E{
				E: 5,
			},
			F: 2.0,
			G: "abc",
		},
		IntPointer: ptr.Int(0),
		List:       []string{"gds", "123"},
		Timeout:    5 * time.Second,
	}, config)
}

func TestUnmarshalWithValidate(t *testing.T) {
	tt := assert.New(t)

	type Struct struct {
		StringRange string           `validate:"@string[10,)"`
		Hostname    string           `validate:"@hostname"`
		Password    presets.Password `validate:"@string[1,)"`
		Struct      struct {
			Int32 int32 `validate:"@int32[2,)"`
		}
	}

	config := &Struct{}
	config.StringRange = "string"
	config.Struct.Int32 = 1

	ok, errMsgs := conf.NewScanner("TEST").Unmarshal(reflect.ValueOf(config), reflect.TypeOf(config))
	tt.False(ok)

	if !ok {
		t.Log(errMsgs)
	}

	assert.Equal(t, transform.ErrMsgMap{
		"Hostname":     "Hostname 类型值错误",
		"Password":     "字符串长度不在[1， 1024]范围内，当前长度：0",
		"StringRange":  "字符串长度不在[10， 1024]范围内，当前长度：6",
		"Struct.Int32": "整形值不在[2, 2147483647)范围内，当前值：1",
	}, errMsgs)
}
