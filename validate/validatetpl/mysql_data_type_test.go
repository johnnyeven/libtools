package validatetpl

import (
	"testing"
)

func TestValidateMySQLDataType(t *testing.T) {
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
			name: "tinyint",
			args: args{
				v: "tinyint(8)",
			},
			want:  true,
			want1: "",
		},
		{
			name: "tinyint unsigned",
			args: args{
				v: "tinyint(8) unsigned",
			},
			want:  true,
			want1: "",
		},
		{
			name: "smallint",
			args: args{
				v: "smallint(16)",
			},
			want:  true,
			want1: "",
		},
		{
			name: "smallint unsigned",
			args: args{
				v: "smallint(16) unsigned",
			},
			want:  true,
			want1: "",
		},
		{
			name: "int",
			args: args{
				v: "int(32)",
			},
			want:  true,
			want1: "",
		},
		{
			name: "int unsigned",
			args: args{
				v: "int(32) unsigned",
			},
			want:  true,
			want1: "",
		},
		{
			name: "bigint",
			args: args{
				v: "bigint(64)",
			},
			want:  true,
			want1: "",
		},
		{
			name: "bigint unsigned",
			args: args{
				v: "bigint(64) unsigned",
			},
			want:  true,
			want1: "",
		},
		{
			name: "varchar",
			args: args{
				v: "varchar(64)",
			},
			want:  true,
			want1: "",
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateMySQLDataType(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateMySQLDataType() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateMySQLDataType() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
