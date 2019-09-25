package amount

import (
	"testing"
)

func TestFenToYuan(t *testing.T) {

	type arg struct {
		input int64
		want  float64
	}

	testCases := []arg{
		{
			input: 1,
			want:  0.01,
		},
		{
			input: 0,
			want:  0.00,
		},
		{
			input: 100,
			want:  1.00,
		},
		{
			input: 101,
			want:  1.01,
		},
		{
			input: 105,
			want:  1.05,
		},
		{
			input: 58805,
			want:  588.05,
		}, {
			input: 50005,
			want:  500.05,
		},
	}

	for _, tc := range testCases {
		output := FenToYuan(tc.input)
		if output != tc.want {
			t.Errorf("input:%d want:%f output:%f", tc.input, tc.want, output)
		} else {
			t.Logf("input:%d want:%f output:%f", tc.input, tc.want, output)
		}
	}
}

func TestYuanToFen(t *testing.T) {
	type arg struct {
		input float64
		want  int64
	}

	testCases := []arg{
		{
			input: 0.01,
			want:  1,
		},
		{
			input: 0.00,
			want:  0,
		},
		{
			input: 1.00,
			want:  100,
		},
		{
			input: 1.01,
			want:  101,
		},
		{
			input: 1.05,
			want:  105,
		},
		{
			input: 588.05,
			want:  58805,
		}, {
			input: 500.05,
			want:  50005,
		}, {
			input: 500.55,
			want:  50055,
		}, {
			input: 500.54,
			want:  50054,
		}, {
			input: 500.54,
			want:  50054,
		}, {
			input: 4.19,
			want:  419,
		},
	}

	for _, tc := range testCases {
		output := YuanToFen(tc.input)
		if output != tc.want {
			t.Errorf("input:%f want:%d output:%d", tc.input, tc.want, output)
		} else {
			t.Logf("input:%f want:%d output:%d", tc.input, tc.want, output)
		}
	}
}

func TestFloat64ToInt64(t *testing.T) {
	type arg struct {
		input   float64
		decimal int
		want    int64
	}

	testCases := []arg{
		{
			input:   1.00,
			decimal: 1,
			want:    10,
		},
		{
			input:   0.01,
			decimal: 2,
			want:    1,
		},
		{
			input:   1.01,
			decimal: 3,
			want:    1010,
		},
		{
			input:   1.05,
			decimal: 4,
			want:    10500,
		},
		{
			input:   588.05,
			decimal: 5,
			want:    58805000,
		}, {
			input:   500.05,
			decimal: 6,
			want:    500050000,
		}, {
			input:   500.55,
			decimal: 7,
			want:    5005500000,
		}, {
			input:   500.54,
			decimal: 8,
			want:    50054000000,
		}, {
			input:   500.54,
			decimal: 9,
			want:    500540000000,
		}, {
			input:   4.19,
			decimal: 10,
			want:    41900000000,
		},
	}

	for _, tc := range testCases {
		output := Float64ToInt64(tc.input, tc.decimal)
		if output != tc.want {
			t.Errorf("input:%f want:%d output:%d", tc.input, tc.want, output)
		} else {
			t.Logf("input:%f want:%d output:%d", tc.input, tc.want, output)
		}
	}
}

func TestInt64ToFloat64(t *testing.T) {

	type arg struct {
		input   int64
		decimal int
		want    float64
	}

	testCases := []arg{
		{
			input:   1,
			decimal: 1,
			want:    0.1,
		},
		{
			input:   0,
			decimal: 2,
			want:    0.00,
		},
		{
			input:   100,
			decimal: 3,
			want:    0.100,
		},
		{
			input:   101,
			decimal: 4,
			want:    0.0101,
		},
		{
			input:   105,
			decimal: 5,
			want:    0.00105,
		},
		{
			input:   58805,
			decimal: 6,
			want:    0.058805,
		}, {
			input:   50005,
			decimal: 7,
			want:    0.0050005,
		},
	}

	for _, tc := range testCases {
		output := Int64ToFloat64(tc.input, tc.decimal)
		if output != tc.want {
			t.Errorf("input:%d want:%f output:%f", tc.input, tc.want, output)
		} else {
			t.Logf("input:%d want:%f output:%f", tc.input, tc.want, output)
		}
	}
}

func TestRound(t *testing.T) {

	type arg struct {
		input    float64
		want     float64
		decimals int
	}
	testCases := []arg{
		{
			input:    1.21,
			want:     1.2,
			decimals: 1,
		},
		{
			input:    1.24,
			want:     1.2,
			decimals: 1,
		},
		{
			input:    1.25,
			want:     1.3,
			decimals: 1,
		},
		{
			input:    1.29,
			want:     1.3,
			decimals: 1,
		},
		{
			input:    1.3,
			want:     1.0,
			decimals: 0,
		},
		{
			input:    1.5,
			want:     2.0,
			decimals: 0,
		},
		{
			input:    1.9,
			want:     2.0,
			decimals: 0,
		},
	}
	for _, tc := range testCases {
		output := Round(tc.input, tc.decimals)
		if output != tc.want {
			t.Errorf("input:%f want:%f output:%f", tc.input, tc.want, output)
		} else {
			t.Logf("input:%f want:%f output:%f", tc.input, tc.want, output)
		}
	}
}

func TestCountInterestByAnnualizedRate(t *testing.T) {

	type arg struct {
		capital  float64
		rate     float64
		decimals int
		dateDiff int
		want     float64
	}
	testCases := []arg{
		{
			capital:  100000,
			rate:     0.05,
			decimals: 2,
			dateDiff: 1,
			want:     13.89,
		},
	}
	for _, tc := range testCases {
		output := CountInterestByAnnualizedRate(tc.capital, tc.rate, tc.dateDiff, tc.decimals)
		if output != tc.want {
			t.Errorf("%+v output:%f", tc, output)
		} else {
			t.Logf("%+v output:%f", tc, output)
		}
	}
}
