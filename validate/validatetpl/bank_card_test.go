package validatetpl

import (
	"testing"
)

func TestValidateBankCard(t *testing.T) {
	type args struct {
		v    interface{}
		ok   bool
		name string
	}

	testCases := []args{
		{
			name: "err",
			v:    "31qadads",
			ok:   false,
		},
		{
			name: "length err",
			v:    "12345678900",
			ok:   false,
		},
		{
			name: "length min",
			v:    "123456789121",
			ok:   true,
		},
		{
			name: "length max",
			v:    "1234567891234567891",
			ok:   true,
		},
		{
			name: "exec length max",
			v:    "12345678912345678912",
			ok:   false,
		},
	}

	for _, tc := range testCases {
		if ok, _ := ValidateBankCard(tc.v); ok != tc.ok {
			t.Errorf("ValidateBankCard input:%s name:%s got:%v want:%v", tc.v, tc.name, ok, tc.ok)
		}
	}
}
