package vincode

import (
	"fmt"
	"testing"
	"unicode"
)

func TestCheckVINLength1(t *testing.T) {
	bytes := make([]byte, 17, 17)
	str := string(bytes)
	expectedLen := vinLegalLength
	actualLen := len(str)
	if expectedLen != actualLen {
		t.Errorf("expected length:%d,actual length:%d", expectedLen, actualLen)
	}
	if !checkVINLength(str) {
		t.Errorf("str:%s length is not correct", str)
	}
}

func TestCheckVINLength2(t *testing.T) {
	bytes := make([]byte, 16, 16)
	str := string(bytes)
	if checkVINLength(str) {
		t.Errorf("str:%s length is not correct", str)
	}
}

func TestCheckVINCharacter1(t *testing.T) {
	bytes := make([]byte, 17, 17)
	for k := range bytes {
		bytes[k] = 'A'
	}
	str := string(bytes)
	if !checkVINCharacter(str) {
		t.Errorf("str:%s has illegal character", str)
	}
}

func TestCheckVINCharacter2(t *testing.T) {
	bytes := make([]byte, 17, 17)
	for k := range bytes {
		bytes[k] = '2'
	}
	str := string(bytes)
	if !checkVINCharacter(str) {
		t.Errorf("str:%s has illegal character", str)
	}
}

func TestCheckVINCharacter3(t *testing.T) {
	str := "中文的VIN1234"

	var tmp rune
	tmp = rune(str[0])
	fmt.Println(unicode.IsLetter(tmp))
	if checkVINCharacter(str) {
		t.Errorf("str:%s has illegal character", str)
	}
}

func TestCheckVINCheckDigit1(t *testing.T) {
	vin := "LFWADRJF011002346"
	if !checkVINCheckDigit(vin) {
		t.Errorf("vin:%s check digit is wrong", vin)
	}
}

func TestCheckVINCheckDigit2(t *testing.T) {
	vin := "LFWADRJF011002356"
	if checkVINCheckDigit(vin) {
		t.Errorf("vin:%s check digit is wrong", vin)
	}
}

func TestVINValidator1(t *testing.T) {
	//for right code
	vinList := []string{"LSVHJ133022221761", "LFWADRJF011002346"}
	for _, v := range vinList {
		if err := VINValidator(v, true); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestVINValidator2(t *testing.T) {
	vin := "LFWADRJF01100346"
	if err := VINValidator(vin, true); err != nil {
		t.Log(err.Error())
	}
}

func TestVINValidator3(t *testing.T) {
	vin := "JA32A3HU6HH123456"
	if err := VINValidator(vin, true); err != nil {
		t.Error(err.Error())
	}
}
