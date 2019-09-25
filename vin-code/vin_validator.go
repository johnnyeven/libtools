package vincode

import (
	"fmt"
	"strings"
)

// VIN 约定
const (
	// VIN 约定长度
	vinLegalLength = 17
)

// 字符对应值 用于计算校验和
var vinLettersValue = map[byte]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 1, 'B': 2, 'C': 3, 'D': 4, 'E': 5, 'F': 6, 'G': 7, 'H': 8,
	'J': 1, 'K': 2, 'L': 3, 'M': 4, 'N': 5,
	'P': 7,
	'R': 9,
	'S': 2, 'T': 3, 'U': 4, 'V': 5, 'W': 6, 'X': 7, 'Y': 8, 'Z': 9,
}

// 位置加权系数
var vinIndexValue = map[int]int{
	0: 8, 1: 7, 2: 6, 3: 5, 4: 4, 5: 3, 6: 2,
	7: 10,
	9: 9, 10: 8, 11: 7, 12: 6, 13: 5, 14: 4, 15: 3, 16: 2,
}

// 统一转换VIN为大写字符
func vinToUpper(vin string) string {
	return strings.ToUpper(vin)
}

// 校验VIN 长度
func checkVINLength(vin string) bool {
	if len(vin) != vinLegalLength {
		return false
	}
	return true
}

// 是否为ASCII大写字母
func isASCIIUpper(n rune) bool {
	return n >= 'A' && n <= 'Z'
}

// 是否为ASCII小写字母
func isASCIILower(n rune) bool {
	return n >= 'a' && n <= 'z'
}

// 是否为ASCII数字
func isASCIINumber(n rune) bool {
	return n >= '0' && n <= '9'
}

// 校验VIN 是否由字母和数字字符构成
func checkVINCharacter(vin string) bool {
	for _, v := range []rune(vin) {
		if !isASCIIUpper(v) && !isASCIILower(v) && !isASCIINumber(v) {
			return false
		}
	}
	return true
}

// 校验VIN 校验位
func checkVINCheckDigit(vin string) bool {
	actualNumRune := '-'
	checkSum := 0
	for k, v := range []byte(vin) {
		if k == 8 {
			actualNumRune = rune(v)
		} else {
			lv, ok := vinLettersValue[v]
			if !ok {
				fmt.Printf("v:%v not in vinLettersValue:%v", v, vinLettersValue)
				return false
			}
			iv, ok := vinIndexValue[k]
			if !ok {
				fmt.Printf("index:%d not in vinIndexValue:%v", k, vinIndexValue)
				return false
			}
			// each character's check num is calc by letter's value multiply index's value
			checkSum += lv * iv
		}
	}

	checkMod := checkSum % 11
	checkDigit := ' '
	if checkMod == 10 {
		checkDigit = 'X'
	} else {
		checkDigit = '0' + rune(checkMod)
	}
	//fmt.Println(checkSum, string(checkDigit))
	if checkDigit != actualNumRune {
		fmt.Printf("expected check digit:%c,actual check digit:%c", checkDigit, actualNumRune)
		return false
	}

	return true
}

func VINValidator(vin string, checkDigit bool) error {
	vinStr := vinToUpper(vin)
	if !checkVINLength(vinStr) {
		return VINCodeLengthError
	}
	if !checkVINCharacter(vinStr) {
		return VINCodeCharacterError
	}
	// for special cars,like volvo has no check digit
	if checkDigit {
		if !checkVINCheckDigit(vinStr) {
			return VINCodeCheckDigitError
		}
	}

	return nil
}
