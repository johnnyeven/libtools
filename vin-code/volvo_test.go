package vincode

import (
	"testing"
)

func TestVolvo1(t *testing.T) {
	vin := "YV1LW65F1Y2123456"
	data, err := VINUnmarshal(vin)
	if err != nil {
		t.Error(vin, err)
	}
	t.Logf("vin:%s,vin data:%+v", vin, data)
}
