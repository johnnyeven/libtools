package validatetpl

import (
	"testing"
)

func TestValidateIPv4(t *testing.T) {
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
				v: "134.255.255.1",
			},
			want:  true,
			want1: "",
		},
		{
			name: "5part",
			args: args{
				v: "134.255.255.1.1",
			},
			want:  false,
			want1: InvalidIPv4Value,
		},
		{
			name: "withLetter",
			args: args{
				v: "134.255.255.2a",
			},
			want:  false,
			want1: InvalidIPv4Value,
		},
		{
			name: "over255",
			args: args{
				v: "134.255.255.256",
			},
			want:  false,
			want1: InvalidIPv4Value,
		},
		{
			name: "intValue",
			args: args{
				v: 255,
			},
			want:  false,
			want1: InvalidIPv4Type,
		},
		{
			name: "invalidIP",
			args: args{
				v: "10.",
			},
			want:  false,
			want1: InvalidIPv4Value,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateIPv4(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateIPv4() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateIPv4() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
