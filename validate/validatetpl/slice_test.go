package validatetpl

import (
	"fmt"
	"testing"
)

func TestValidateSlice(t *testing.T) {
	type args struct {
		min     uint64
		max     uint64
		elemMin uint64
		elemMax uint64
		v       interface{}
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
				min:     0,
				max:     0,
				elemMin: 0,
				elemMax: 10,
				v: []uint64{
					1,
					2,
					3,
					4,
				},
			},
			want:  true,
			want1: "",
		},
		{
			name: "notSLice",
			args: args{
				min:     0,
				max:     0,
				elemMin: 0,
				elemMax: 10,
				v:       1,
			},
			want:  false,
			want1: TYPE_NOT_SLICE,
		},
		{
			name: "elemLtRequire",
			args: args{
				min:     5,
				max:     10,
				elemMin: 0,
				elemMax: 10,
				v: []uint64{
					1,
					2,
					3,
					4,
				},
			},
			want:  false,
			want1: fmt.Sprintf(SLICE_ELEM_NUM_NOT_IN_RANGE, 5, 10, 4),
		},
		{
			name: "elemGtRequire",
			args: args{
				min:     1,
				max:     3,
				elemMin: 0,
				elemMax: 10,
				v: []uint64{
					1,
					2,
					3,
					4,
				},
			},
			want:  false,
			want1: fmt.Sprintf(SLICE_ELEM_NUM_NOT_IN_RANGE, 1, 3, 4),
		},
		{
			name: "elemValueTypeInvalid",
			args: args{
				min:     0,
				max:     0,
				elemMin: 0,
				elemMax: 10,
				v: []int64{
					1,
					2,
					3,
					4,
				},
			},
			want:  false,
			want1: fmt.Sprintf(SLICE_ELEM_INVALID, TYPE_NOT_UINT64),
		},
		{
			name: "elemValueNotInRange",
			args: args{
				min:     0,
				max:     0,
				elemMin: 0,
				elemMax: 10,
				v: []uint64{
					11,
				},
			},
			want:  false,
			want1: fmt.Sprintf(SLICE_ELEM_INVALID, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, "[", 0, 10, "]", 11)),
		},
	}
	for _, tt := range tests {
		itemValidateFunc := NewRangeValidateUint64(tt.args.elemMin, tt.args.elemMax, false, false)
		sliceValidateFunc := NewValidateSlice(tt.args.min, tt.args.max, itemValidateFunc)
		got, got1 := sliceValidateFunc(tt.args.v)
		if got != tt.want {
			t.Errorf("%q. ValidateSlice() got = %v, want %v", tt.name, got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("%q. ValidateSlice() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
