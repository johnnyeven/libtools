package misc

import (
	"testing"
)

func TestGetMFRS(t *testing.T) {
	vin1 := "JA45"
	vin1Expect := "Mitsubishi"
	vin2 := "JA23"
	vin2Expect := "Isuzu"
	vin3 := "AAAA"
	vin3Expect := UnknownMFRS

	if vin1Expect != GetVINManuf(vin1) {
		t.Error(vin1, vin1Expect, GetVINManuf(vin1))
	}
	if vin2Expect != GetVINManuf(vin2) {
		t.Error(vin2, vin2Expect, GetVINManuf(vin2))
	}
	if vin3Expect != GetVINManuf(vin3) {
		t.Error(vin3, vin3Expect, GetVINManuf(vin3))
	}

}
