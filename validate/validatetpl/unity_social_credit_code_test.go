package validatetpl

import (
	"testing"
)

func TestValidateUnitySocialCreditCode(t *testing.T) {
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
				v: "12345678901234567a",
			},
			want:  true,
			want1: "",
		},
		{
			name: "ExceedLength",
			args: args{
				v: "1234567890123456789",
			},
			want:  false,
			want1: InvalidUnitySocialCreditCodeValue,
		},
		{
			name: "shortLength",
			args: args{
				v: "12345678901234567",
			},
			want:  false,
			want1: InvalidUnitySocialCreditCodeValue,
		},
		{
			name: "invalidCharacter",
			args: args{
				v: "12345678901234567@",
			},
			want:  false,
			want1: InvalidUnitySocialCreditCodeValue,
		},
		{
			name: "invalidType",
			args: args{
				v: 123456789012345678,
			},
			want:  false,
			want1: InvalidUnitySocialCreditCodeType,
		},
		{
			name: "empty",
			args: args{
				v: "",
			},
			want:  false,
			want1: InvalidUnitySocialCreditCodeValue,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateUnitySocialCreditCode(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateUnitySocialCreditCode() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateUnitySocialCreditCode() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
