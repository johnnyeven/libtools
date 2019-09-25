package validatetpl

import (
	"testing"
)

func TestValidateUrl(t *testing.T) {
	too_long_url := ""
	for i := 0; i < MAX_URL_LEN+1; i++ {
		too_long_url += "1"
	}

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
				v: "http://user:passwd@127.0.0.1:8080/abdf?q=a&dks",
			},
			want:  true,
			want1: "",
		},
		{
			name: "normal1",
			args: args{
				v: "https://user:passwd@127.0.0.1:8080/abdf?q=a&dks",
			},
			want:  true,
			want1: "",
		},
		{
			name: "error",
			args: args{
				v: "//userpasswd@127.0.0.1:8080/abdf?q=a&dks",
			},
			want:  false,
			want1: HTTP_URL_SCHEME_ERROR,
		},
		{
			name: "too long",
			args: args{
				v: too_long_url,
			},
			want:  false,
			want1: HTTP_URL_TOO_LONG,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateHttpUrl(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateUrl() got = %v, want %v", tt.name, got, tt.want)
		}

		if got1 != tt.want1 {
			t.Errorf("%q. ValidateUrl() got = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}

func TestValidateHttpUrlOrEmpty(t *testing.T) {
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
				v: "http://user:passwd@127.0.0.1:8080/abdf?q=a&dks",
			},
			want:  true,
			want1: "",
		},
		{
			name: "emptyNormal",
			args: args{
				v: "",
			},
			want:  true,
			want1: "",
		},
		{
			name: "error",
			args: args{
				v: "//userpasswd@127.0.0.1:8080/abdf?q=a&dks",
			},
			want:  false,
			want1: HTTP_URL_SCHEME_ERROR,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateHttpUrlOrEmpty(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateHttpUrlOrEmpty() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateHttpUrlOrEmpty() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
