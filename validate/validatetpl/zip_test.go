package validatetpl

import (
	"testing"
)

func TestValidateZip(t *testing.T) {
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
				v: "123456",
			},
			want:  true,
			want1: "",
		},
		{
			name: "alpha",
			args: args{
				v: "123a56",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
		{
			name: "more",
			args: args{
				v: "1234568",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
		{
			name: "more2",
			args: args{
				v: "123456a",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
		{
			name: "less",
			args: args{
				v: "asdf",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
	}

	for _, tt := range tests {
		got, got1 := ValidateZipZh(tt.args.v)
		t.Log(got, got1)
		if got != tt.want {
			t.Errorf("%q. ValidateZipZh() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateZipZh() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}

func TestValidateZipOrEmpty(t *testing.T) {
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
				v: "123456",
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
		{
			name: "more",
			args: args{
				v: "1234568",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
		{
			name: "more2",
			args: args{
				v: "123456a",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
		{
			name: "less",
			args: args{
				v: "asdf",
			},
			want:  false,
			want1: InvalidZipZhValue,
		},
	}

	for _, tt := range tests {
		got, got1 := ValidateZipZhOrEmpty(tt.args.v)
		t.Log(got, got1)
		if got != tt.want {
			t.Errorf("%q. ValidateZipZhOrEmpty() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateZipZhOrEmpty() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
