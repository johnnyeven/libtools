package vincode

import (
	"strings"

	"github.com/johnnyeven/libtools/vin-code/misc"
)

// 专业术语(Technical Terms)
// VIN = "车辆识别代号(vehicle identification number)"
// WMI = "世界制造厂识别代号(world manufacturer identifier)"
// VDS = "车辆说明部分(vehicle descriptor section)"
// VIS = "车辆指示部分(vehicle indicator section)"

func vinUnmarshal(vin string) (VINData, error) {
	//vinData := VINData{}
	//country := misc.GetVINCountry(vin)
	mfrs := misc.GetVINManuf(vin)
	switch {
	case strings.Contains(mfrs, "Mitsubishi"):
		// 三菱2017
		msbsVIN := MSBSVINCode(vin)
		return ParseAll(msbsVIN)
	case strings.Contains(mfrs, "Volvo"):
		// 沃尔沃
		volVIN := VOLVOVINCode(vin)
		return ParseAll(volVIN)
	default:
		// check digit
		if err := VINValidator(vin, true); err != nil {
			vinData := VINData{}
			return vinData, err
		}
		generalVIN := GeneralVINCode(vin)
		return ParseAll(generalVIN)
	}
}

func VINUnmarshal(vinStr string) (VINData, error) {
	vin := strings.ToUpper(vinStr)
	if err := VINValidator(vin, false); err != nil {
		vinData := VINData{}
		return vinData, err
	}
	return vinUnmarshal(vin)
}
