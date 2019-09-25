package validatetpl

import (
	"testing"
)

func TestValidateCrontab(t *testing.T) {
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
			name: "per2sec",
			args: args{
				v: "*/2 * * * * ?",
			},
			want:  true,
			want1: "",
		},
		{
			name: "empty",
			args: args{
				v: "",
			},
			want:  false,
			want1: InvalidCrontabValue,
		},
		{
			name: "wrong type",
			args: args{
				v: 12,
			},
			want:  false,
			want1: InvalidCrontabType,
		},
		{
			name: "wrong contents",
			args: args{
				v: "112dsww",
			},
			want:  false,
			want1: InvalidCrontabValue,
		},
	}
	for _, tt := range tests {
		got, got1 := ValidateCrontab(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateCrontab() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateCrontab() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
