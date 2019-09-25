package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnnyeven/libtools/validate/validatetpl"
)

func TestValidateItemWithPredefineTag(t *testing.T) {
	char := string("中文中文")
	if ok, str := ValidateItem("@char[0,2]", char); ok {
		t.Log(str)
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_uint8 := uint8(0)
	if ok, str := ValidateItem("@uint8[1,]", value_uint8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_int8 := int8(-100)
	if ok, str := ValidateItem("@int8[-99,-10]", value_int8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}

	value_uint16 := uint16(0)
	if ok, str := ValidateItem("@uint16[1,]", value_uint16); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_int16 := int16(-100)
	if ok, str := ValidateItem("@int16[-99,0]", value_int16); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_uint32 := uint32(0)
	if ok, str := ValidateItem("@uint32[1,]", value_uint32); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_int32 := int32(-100)
	if ok, str := ValidateItem("@int32[-99,0]", value_int32); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_uint64 := uint64(0)
	if ok, str := ValidateItem("@uint64[1,]", value_uint64); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_int64 := int64(-100)
	if ok, str := ValidateItem("@int64[-99,1]", value_int64); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_string := string("hellow word")
	if ok, str := ValidateItem("@string[0,10]", value_string); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_unlimit_string := string("unlimit")
	if ok, str := ValidateItem("@string[0,unlimit]", value_unlimit_string); !ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	value_mobile := string("28520270839")
	if ok, str := ValidateItem("@regexp[^((\\+86)|(86))?1\\d{10}$]", value_mobile); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	// test [,]
	value_int8 = int8(1)
	if ok, str := ValidateItem("@int8[,1]", value_int8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_int16 = int16(1)
	if ok, str := ValidateItem("@int16[,1]", value_int16); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_int32 = int32(1)
	if ok, str := ValidateItem("@int32[,1]", value_int32); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_int64 = int64(1)
	if ok, str := ValidateItem("@int64[,1]", value_int64); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}

	value_uint8 = uint8(1)
	if ok, str := ValidateItem("@uint8[,1]", value_uint8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_uint16 = uint16(1)
	if ok, str := ValidateItem("@uint16[,1]", value_uint16); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_uint32 = uint32(1)
	if ok, str := ValidateItem("@uint32[,1]", value_uint32); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_uint64 = uint64(0)
	if ok, str := ValidateItem("@uint64[,1]", value_uint64); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_float32 := float32(1.1)
	if ok, str := ValidateItem("@float32[1.1,1.2]", value_float32); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	value_float64 := float64(0)
	if ok, str := ValidateItem("@float64[,1.2]", value_float64); !ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}

	value_ip := string("1.0.1")
	if ok, str := ValidateItem("@ipv4", value_ip); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
}

func TestValidateInt8ItemWithPredefineTag(t *testing.T) {
	// left [
	value_int8 := int8(10)
	if ok, str := ValidateItem("@int8[10,20]", value_int8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	// right ]
	if ok, str := ValidateItem("@int8[,10]", value_int8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	// left (
	if ok, str := ValidateItem("@int8(10,20]", value_int8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	// right )
	if ok, str := ValidateItem("@int8[,10)", value_int8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	// enum
	if ok, str := ValidateItem("@int8{1,2,3}", value_int8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	if ok, str := ValidateItem("@int8{1,2,3,10}", value_int8); !ok {
		t.Error("should be valid")
		t.Log(str)
	} else {
		t.Log(str)
	}
}

func TestValidateUint8ItemWithPredefineTag(t *testing.T) {
	// left [
	value_uint8 := uint8(10)
	if ok, str := ValidateItem("@uint8[10,20]", value_uint8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	// right ]
	if ok, str := ValidateItem("@uint8[,10]", value_uint8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}
	// left (
	if ok, str := ValidateItem("@uint8(10,20]", value_uint8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	// right )
	if ok, str := ValidateItem("@uint8[,10)", value_uint8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
}

func TestValidateEnumWithPredefineTag(t *testing.T) {
	value_uint8 := uint8(10)
	// enum
	if ok, str := ValidateItem("@uint8{1,2,3}", value_uint8); ok {
		t.Error("should be invalid")
	} else {
		t.Log(str)
	}
	if ok, str := ValidateItem("@uint8{1,2,3,10}", value_uint8); !ok {
		t.Error("should be valid")
	} else {
		t.Log(str)
	}

	var tStr string = "NORMAL"
	if ok, str := ValidateItem("@string{NORMAL}", tStr); !ok {
		t.Error("should be valid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// with blank
	if ok, str := ValidateItem("@string{,NORMAL}", tStr); !ok {
		t.Error("should be valid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// with _
	if ok, str := ValidateItem("@string{NORMAL_A}", tStr); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}
}

func TestValidateSliceWithPredefineTag(t *testing.T) {
	// normal
	var slice []int32 = []int32{1, 2, 3}
	if ok, str := ValidateItem("@array[1,3]:@int32[1,3]", slice); !ok {
		t.Error("should be valid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// elem type invalid
	var invalidSlice []uint32 = []uint32{1, 2, 3}
	if ok, str := ValidateItem("@array[1,3]:@int32[1,3]", invalidSlice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// elem too few
	slice = []int32{}
	if ok, str := ValidateItem("@array[1,3]:@int32[1,3]", slice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// elem too many
	slice = []int32{1, 2, 2, 3}
	if ok, str := ValidateItem("@array[1,3]:@int32[1,3]", slice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// elem too small
	slice = []int32{0, 2, 3}
	if ok, str := ValidateItem("@array[1,3]:@int32[1,3]", slice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}
	// elem too big
	slice = []int32{1, 2, 4}
	if ok, str := ValidateItem("@array[1,3]:@int32[1,3]", slice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}

	AddValidateFunc("@ipv4", validatetpl.ValidateIPv4)
	strSlice := []string{"1.0"}
	if ok, str := ValidateItem("@array[0,]:@ipv4", strSlice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}

	// invalid item validate
	slice = []int32{1, 2, 4, 4}
	if ok, str := ValidateItem("@array[1,3]:@invalid", slice); ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}

	// invalid item validate
	slice = []int32{1, 2, 4, 4}
	if ok, str := ValidateItem("@array[1,4]", slice); !ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}

	// invalid item validate
	slice = []int32{1, 2, 4, 4}
	if ok, str := ValidateItem("@array[1,4]:@", slice); !ok {
		t.Error("should be invalid")
		t.Log(str)
	} else {
		t.Log(str)
	}
}

func TestFloat32Range(t *testing.T) {
	valueFloat32 := float32(99.9999)
	ok, str := ValidateItem("@float32[0,99]", valueFloat32)
	assert.False(t, ok)
	assert.Equal(t, "浮点值不在[0,99]范围内，当前值：99.9999", str)
}

func TestFloat64Range(t *testing.T) {
	value_float64 := float64(321.1234567890123456789)
	ok, str := ValidateItem("@float64[1,1.2]", value_float64)
	assert.False(t, ok)
	assert.Equal(t, "浮点值不在[1,1.2]范围内，当前值：321.1234567890123", str)
}

func TestFloat32Decimal(t *testing.T) {
	value_float32 := float32(0.12345)
	ok, str := ValidateItem("@float32<7,4>[0,99]", value_float32)
	assert.False(t, ok)
	assert.Equal(t, "浮点值小数位必须为4，总位数不能超过7位，当前值：0.12345", str)
}

func TestFloat64Decimal(t *testing.T) {
	value_float64 := float64(99999.999)
	ok, str := ValidateItem("@float64<8,4>", value_float64)
	assert.False(t, ok)
	assert.Equal(t, "浮点值小数位必须为4，总位数不能超过8位，当前值：99999.999", str)
}
