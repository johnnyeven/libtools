package msbs

import (
	"fmt"

	"github.com/johnnyeven/libtools/vin-code/vinrune"
)

var (
	// 制造国 pos 1
	Country = map[rune]string{
		'J': "Japan (MMC)",
		'M': "Thailand (MMT)",
	}
	// 制造商 pos 2
	Manufact = map[rune]string{
		'A': "Mitsubishi",
		'L': "Mitsubishi Motors Thailand",
	}
	// 车用类型 pos 3
	VehicleType = map[rune]string{
		'3': "Passenger Car",
		'4': "Multi−Purpose Vehicle",
	}
	// 约束系统 pos 4
	RestraintSys = map[rune]string{
		'2': "1st & 2nd Row Curtain + Seat Air Bags; Passenger Car",
		'A': "1st & 2nd Row Curtain + 1st Row Seat Air Bags; MPV up to 5,000 lbs GVWR",
		'J': "1st & 2nd Row Curtain + 1st Row Seat Air Bags; MPV over 5,000 lbs GVWR",
	}
	// 车系 pos 5 & 6
	CarSeries = map[string]string{
		"15": "Mitsubishi i−MiEV ES",
		"A3": "Mitsubishi Mirage ES",
		"A4": "Mitsubishi Mirage SE",
		"A5": "Mitsubishi Mirage GT (SEL in Canada)",
		"F3": "Mitsubishi Mirage G4 ES",
		"F4": "Mitsubishi Mirage G4 SE (SEL in Canada)",
		"D2": "Mitsubishi Outlander ES (FWD)",
		"D3": "Mitsubishi Outlander SE/SEL (FWD)",
		"Z2": "Mitsubishi Outlander ES (AWC)",
		"Z3": "Mitsubishi Outlander SE/SEL (S−AWC)",
		//"Z3": "Mitsubishi Outlander SE (AWC) (Canada only)",
		"Z4": "Mitsubishi Outlander GT (S−AWC)",
		"U2": "Mitsubishi Lancer ES (FWD)",
		"V2": "Mitsubishi Lancer ES/SE/SEL (AWC)",
		//"U2": "Mitsubishi Lancer SE/SE Limited (FWD) (Canada only)",
		//"V2": "Mitsubishi Lancer SE Limited (AWC) (Canada only)",
		//"V2": "Mitsubishi Lancer GTS (AWC) (Canada only)",
		//"U8": "Mitsubishi Lancer GTS (FWD) (Canada only)",
		//"X2": "Mitsubishi Lancer Sportback SE Limited (Canada only)",
		//"X2": "Mitsubishi Lancer Sportback GT (Canada only)",
		//"H3": "Mitsubishi RVR ES/SE (FWD) (Canada only)",
		//"J3": "Mitsubishi RVR SE (AWC) (Canada only)",
		//"J4": "Mitsubishi RVR SE Limited/GT (AWC) (Canada only)",
		"P3": "Mitsubishi Outlander Sport ES/SE (FWD)",
		"P4": "Mitsubishi Outlander Sport SEL (FWD)",
		"R3": "Mitsubishi Outlander Sport ES/SE (AWC)",
		"R4": "Mitsubishi Outlander Sport SEL/GT (AWC)",
	}
	// 车外观类型 pos 7
	DoorType = map[rune]string{
		'H': "5−door Hatchback (i−MiEV, Mirage)",
		'F': "4−door Sedan (Mirage G4, Lancer)",
		'A': "5−door Wagon/SUV (Outlander, Outlander Sport/RVR)",
	}

	// 发动机 pos 8
	Engine = map[rune]string{
		'J': "1.2L DOHC MIVEC (3A92) Gasoline",
		'4': "49Kw Electric Motor (Y51)",
		'U': "2.0L DOHC MIVEC (4B11) Gasoline",
		'W': "2.4L DOHC MIVEC (4B12) Gasoline",
		'3': "2.4L MIVEC (4J12) Gasoline",
		'X': "3.0L MIVEC (6B31) Gasoline",
	}
	// 年份 pos 10
	// 装配厂 pos 11
	AssemblyPlant = map[rune]string{
		'H': "Laem Chabang−3 (Thailand)",
		'J': "Nagoya (Japan)",
		'U': "Mizushima (Japan)",
		'Z': "Okazaki (Japan)",
	}
	// 顺序号 pos 12~17
)

func GetWMIRune(vin string) vinrune.WMIRune {
	ret := vinrune.WMIRune{}
	vinStr := vin[:3]
	fmt.Sscanf(vinStr, "%c%c%c", &ret.CountryRune, &ret.ManfRune, &ret.VehicleTypeRune)
	return ret
}

func GetVDSRune(vin string) vinrune.VDSRune {
	ret := vinrune.VDSRune{}
	vinStr := vin[3:8]
	fmt.Sscanf(vinStr, "%c%2s%c%c", &ret.RestraintSysRune, &ret.CarSeriesStr, &ret.DoorTypeRune, &ret.EngineRune)
	return ret
}

func GetVISRune(vin string) vinrune.VISRune {
	ret := vinrune.VISRune{}
	vinStr := vin[9:]
	fmt.Sscanf(vinStr, "%c%c%6s", &ret.YearRune, &ret.AssemblyRune, &ret.SequenceNO)
	return ret
}
