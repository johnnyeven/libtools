package vincode

import (
	"github.com/johnnyeven/libtools/vin-code/mfrs/volvo"
	"github.com/johnnyeven/libtools/vin-code/misc"
)

/**********************************************************
** 沃尔沃 风险:
** 1.国标第9位为校验位而沃尔沃第9位可能是变速箱代码
***********************************************************/

// 沃尔沃
type VOLVOVINCode string

func (str VOLVOVINCode) ParseWMI() (WMIData, error) {
	re := volvo.GetWMIRune(string(str))
	wmi := WMIData{}
	wmi.Continent = misc.GetVINContinent(string(str))
	wmi.Country = misc.GetVINCountry(string(str))
	wmi.Manufacturer = misc.GetVINManuf(string(str))

	tv, ok := volvo.VehicleType[re.VehicleTypeStr]
	if !ok {
		return wmi, VINCodeParseVehicleTypeError
	}
	wmi.VehicleType = tv

	return wmi, nil
}

func (str VOLVOVINCode) ParseVDS() (VDSData, error) {
	re := volvo.GetVDSRune(string(str))
	vds := VDSData{}

	rt, ok := volvo.RestraintSys[re.RestraintSysRune]
	if !ok {
		return vds, VINCodeParseRestriantError
	}
	vds.RestraintSystem = rt

	cs, ok := volvo.CarSeries[re.CarSeriesRune]
	if !ok {
		return vds, VINCodeParseCarSeriesError
	}
	vds.CarSeries = cs

	ei, ok := volvo.EmissisonCode[re.EmissisonRune]
	if !ok {
		// oops, emission data not enough
		//return vds, VINCodeParseEmissionError
	}
	vds.Emissison = ei

	eg, ok := volvo.Engine[re.EngineStr]
	if !ok {
		return vds, VINCodeParseEngineError
	}
	vds.Engine = eg

	return vds, nil
}

func (str VOLVOVINCode) ParseVIS() (VISData, error) {
	re := volvo.GetVISRune(string(str))
	vis := VISData{}
	vis.SequenceNO = re.SequenceNO
	ap, ok := volvo.AssemblyPlant[re.AssemblyRune]
	if !ok {
		return vis, VINCodeParseAssemblyError
	}
	vis.AssemblyPlant = ap
	vis.ModelYear = misc.GetModelYearStr(re.YearRune)
	if vis.ModelYear == "0" {
		return vis, VINCodeParseYearError
	}

	return vis, nil
}
