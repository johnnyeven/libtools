package validatetpl

import (
	"testing"
)

func TestValidateOrgInsitituteCode(t *testing.T) {
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
				v: "12345678-A",
			},
			want:  true,
			want1: "",
		},
		{
			name: "ExceedLength",
			args: args{
				v: "123456789-0",
			},
			want:  false,
			want1: InvalidOrgInsitituteCodeValue,
		},
		{
			name: "shortLength",
			args: args{
				v: "1234567-8",
			},
			want:  false,
			want1: InvalidOrgInsitituteCodeValue,
		},
		{
			name: "invalidCharacter",
			args: args{
				v: "12345678{",
			},
			want:  false,
			want1: InvalidOrgInsitituteCodeValue,
		},
		{
			name: "invalidFormat",
			args: args{
				v: "123456789",
			},
			want:  false,
			want1: InvalidOrgInsitituteCodeValue,
		},
		{
			name: "invalidType",
			args: args{
				v: 123456789,
			},
			want:  false,
			want1: InvalidOrgInsitituteCodeType,
		},
		{
			name: "empty",
			args: args{
				v: "",
			},
			want:  false,
			want1: InvalidOrgInsitituteCodeValue,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateOrgInsitituteCode(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateOrgInsitituteCode() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateOrgInsitituteCode() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
