package validatetpl

import (
	"testing"
)

func TestValidatePhone(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			name: "normal",
			args: args{
				v: "13478902992",
			},
			want:  true,
			want1: "",
		},
		{
			name: "ExceedLength",
			args: args{
				v: "134789029932",
			},
			want:  false,
			want1: InvalidPhoneNoValue,
		},
		{
			name: "shortLength",
			args: args{
				v: "1347890292",
			},
			want:  false,
			want1: InvalidPhoneNoValue,
		},
		{
			name: "notAllNumber",
			args: args{
				v: "134789a2992",
			},
			want:  false,
			want1: InvalidPhoneNoValue,
		},
		{
			name: "intType",
			args: args{
				v: 13478902992,
			},
			want:  false,
			want1: InvalidPhoneNoType,
		},
		{
			name: "empty",
			args: args{
				v: "",
			},
			want:  false,
			want1: InvalidPhoneNoValue,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidatePhone(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidatePhone() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidatePhone() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
