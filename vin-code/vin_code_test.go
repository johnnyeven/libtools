package vincode

import (
	"testing"
)

func TestVINCodeUnmarshal1(t *testing.T) {
	vin := "JA32A3HU6HH123456"
	data, err := VINUnmarshal(vin)
	if err != nil {
		t.Error(vin, err)
	}
	t.Logf("vin:%s,vin data:%+v", vin, data)
}
