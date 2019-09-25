package validatetpl

import (
	"testing"
)

func TestValidatePlateNo(t *testing.T) {
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
				v: "川A0xv12",
			},
			want:  true,
			want1: "",
		},
		{
			name: "新能源车6位",
			args: args{
				v: "川A0xv123",
			},
			want:  true,
			want1: "",
		},
		{
			name: "挂车",
			args: args{
				v: "川A0xv1挂",
			},
			want:  true,
			want1: "",
		},
		{
			name: "empty",
			args: args{
				v: "",
			},
			want:  true,
			want1: "",
		},
	}
	for _, tt := range tests {
		got, got1 := ValidatePlateNo(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidatePlateNo() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidatePlateNo() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
