package strutil_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/johnnyeven/libtools/strutil"
)

type Stringer struct {
}

func (Stringer) String() string {
	return "test"
}

func TestConvertToStr(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				v: "abc",
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name: "int",
			args: args{
				v: int(112),
			},
			want:    "112",
			wantErr: false,
		},
		{
			name: "int8",
			args: args{
				v: int8(12),
			},
			want:    "12",
			wantErr: false,
		},
		{
			name: "int16",
			args: args{
				v: int16(123),
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "int32",
			args: args{
				v: int32(456),
			},
			want:    "456",
			wantErr: false,
		},
		{
			name: "int64",
			args: args{
				v: int64(12390),
			},
			want:    "12390",
			wantErr: false,
		},
		{
			name: "uint",
			args: args{
				v: uint(12390),
			},
			want:    "12390",
			wantErr: false,
		},
		{
			name: "uint8",
			args: args{
				v: uint8(12),
			},
			want:    "12",
			wantErr: false,
		},
		{
			name: "uint16",
			args: args{
				v: uint16(12390),
			},
			want:    "12390",
			wantErr: false,
		},
		{
			name: "uint32",
			args: args{
				v: uint32(12390),
			},
			want:    "12390",
			wantErr: false,
		},
		{
			name: "uint64",
			args: args{
				v: uint64(12390),
			},
			want:    "12390",
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want:    "true",
			wantErr: false,
		},
		{
			name: "stringer",
			args: args{
				v: Stringer{},
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "float32",
			args: args{
				v: float32(1.1),
			},
			want:    "1.1",
			wantErr: false,
		},
		{
			name: "float64",
			args: args{
				v: float64(1.11),
			},
			want:    "1.11",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		got, err := strutil.ConvertToStr(tt.args.v)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. ConvertToStr() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. ConvertToStr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestConvertFromStr(t *testing.T) {
	var int8V int8
	var uint8V uint8
	var int16V int16
	var uint16V uint16
	var stringV string
	var boolV bool
	var float32V float32
	var duration time.Duration
	type args struct {
		strValue string
		v        reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "int8",
			args: args{
				strValue: "1",
				v:        reflect.ValueOf(&int8V),
			},
			wantErr: false,
		},
		{
			name: "int8",
			args: args{
				strValue: "100000000",
				v:        reflect.ValueOf(&int8V),
			},
			wantErr: true,
		},
		{
			name: "int8",
			args: args{
				strValue: "-1",
				v:        reflect.ValueOf(&int8V),
			},
			wantErr: false,
		},
		{
			name: "int8",
			args: args{
				strValue: "-1000000000",
				v:        reflect.ValueOf(&int8V),
			},
			wantErr: true,
		},
		{
			name: "uint8",
			args: args{
				strValue: "1",
				v:        reflect.ValueOf(&uint8V),
			},
			wantErr: false,
		},
		{
			name: "uint8",
			args: args{
				strValue: "-1",
				v:        reflect.ValueOf(&uint8V),
			},
			wantErr: true,
		},
		{
			name: "uint8",
			args: args{
				strValue: "10000000",
				v:        reflect.ValueOf(&uint8V),
			},
			wantErr: true,
		},
		{
			name: "int16",
			args: args{
				strValue: "32768",
				v:        reflect.ValueOf(&int16V),
			},
			wantErr: true,
		},
		{
			name: "int16",
			args: args{
				strValue: "-32769",
				v:        reflect.ValueOf(&int16V),
			},
			wantErr: true,
		},
		{
			name: "int16",
			args: args{
				strValue: "32767",
				v:        reflect.ValueOf(&int16V),
			},
			wantErr: false,
		},
		{
			name: "uint16",
			args: args{
				strValue: "65536",
				v:        reflect.ValueOf(&uint16V),
			},
			wantErr: true,
		},
		{
			name: "uint16",
			args: args{
				strValue: "-1",
				v:        reflect.ValueOf(&uint16V),
			},
			wantErr: true,
		},
		{
			name: "uint16",
			args: args{
				strValue: "65535",
				v:        reflect.ValueOf(&uint16V),
			},
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				strValue: "abd-)123$123",
				v:        reflect.ValueOf(&stringV),
			},
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				strValue: "false",
				v:        reflect.ValueOf(&boolV),
			},
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				strValue: "False",
				v:        reflect.ValueOf(&boolV),
			},
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				strValue: "true",
				v:        reflect.ValueOf(&boolV),
			},
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				strValue: "True",
				v:        reflect.ValueOf(&boolV),
			},
			wantErr: false,
		},
		{
			name: "float32",
			args: args{
				strValue: "0.32",
				v:        reflect.ValueOf(&float32V),
			},
			wantErr: false,
		},
		{
			name: "duration",
			args: args{
				strValue: "30s",
				v:        reflect.ValueOf(&duration),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := strutil.ConvertFromStr(tt.args.strValue, tt.args.v); (err != nil) != tt.wantErr {
			t.Errorf("%q. ConvertFromStr() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
