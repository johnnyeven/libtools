package validatetpl

import (
	"testing"
)

func TestValidateIDCardNo(t *testing.T) {
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
				v: "450821198801190031",
			},
			want:  true,
			want1: "",
		},
		{
			name: "X",
			args: args{
				v: "11010719870304001X",
			},
			want:  true,
			want1: "",
		},
		{
			name: "x",
			args: args{
				v: "11010719870304001x",
			},
			want:  true,
			want1: "",
		},
		{
			name: "invalid",
			args: args{
				v: "450821198801190032",
			},
			want:  false,
			want1: InvalidIDCardNoValue,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateIDCardNo(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateIDCardNo() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateIDCardNo() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
