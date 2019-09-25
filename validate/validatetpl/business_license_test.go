package validatetpl

import (
	"testing"
)

func TestValidateBusinessLicense(t *testing.T) {
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
				v: "123456789012345",
			},
			want:  true,
			want1: "",
		},
		{
			name: "ExceedLength",
			args: args{
				v: "1234567890123456",
			},
			want:  false,
			want1: InvalidBusinessLicenseValue,
		},
		{
			name: "shortLength",
			args: args{
				v: "12345678901234",
			},
			want:  false,
			want1: InvalidBusinessLicenseValue,
		},
		{
			name: "invalidCharacter",
			args: args{
				v: "12345678901234a",
			},
			want:  false,
			want1: InvalidBusinessLicenseValue,
		},
		{
			name: "invalidType",
			args: args{
				v: 123456789012345,
			},
			want:  false,
			want1: InvalidBusinessLicenseType,
		},
		{
			name: "empty",
			args: args{
				v: "",
			},
			want:  false,
			want1: InvalidBusinessLicenseValue,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateBusinessLicense(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateBusinessLicense() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateBusinessLicense() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
