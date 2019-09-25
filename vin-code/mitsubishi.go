package vincode

import (
	"github.com/johnnyeven/libtools/vin-code/mfrs/msbs"
	"github.com/johnnyeven/libtools/vin-code/misc"
)

// 三菱
type MSBSVINCode string

func (str MSBSVINCode) ParseWMI() (WMIData, error) {
	re := msbs.GetWMIRune(string(str))
	wmi := WMIData{}
	cv, ok := msbs.Country[re.CountryRune]
	if !ok {
		return wmi, VINCodeParseCountryError
	}
	wmi.Country = cv

	mv, ok := msbs.Manufact[re.ManfRune]
	if !ok {
		return wmi, VINCodeParseManufError
	}
	wmi.Manufacturer = mv

	tv, ok := msbs.VehicleType[re.VehicleTypeRune]
	if !ok {
		return wmi, VINCodeParseVehicleTypeError
	}
	wmi.VehicleType = tv

	return wmi, nil
}

func (str MSBSVINCode) ParseVDS() (VDSData, error) {
	re := msbs.GetVDSRune(string(str))
	vds := VDSData{}

	rt, ok := msbs.RestraintSys[re.RestraintSysRune]
	if !ok {
		return vds, VINCodeParseRestriantError
	}
	vds.RestraintSystem = rt

	cs, ok := msbs.CarSeries[re.CarSeriesStr]
	if !ok {
		return vds, VINCodeParseCarSeriesError
	}
	vds.CarSeries = cs

	dt, ok := msbs.DoorType[re.DoorTypeRune]
	if !ok {
		return vds, VINCodeParseDoorTypeError
	}
	vds.DoorType = dt

	eg, ok := msbs.Engine[re.EngineRune]
	if !ok {
		return vds, VINCodeParseEngineError
	}
	vds.Engine = eg

	return vds, nil
}

func (str MSBSVINCode) ParseVIS() (VISData, error) {
	re := msbs.GetVISRune(string(str))
	vis := VISData{}
	vis.SequenceNO = re.SequenceNO
	ap, ok := msbs.AssemblyPlant[re.AssemblyRune]
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
