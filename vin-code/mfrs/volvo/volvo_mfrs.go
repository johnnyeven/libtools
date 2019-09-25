package volvo

import (
	"fmt"

	"github.com/johnnyeven/libtools/vin-code/vinrune"
)

var (
	// 车用类型 pos 1-3
	VehicleType = map[string]string{
		"XLB": "passenger cars built by Volvo Car BV / NedCar",
		"YV1": "passenger cars",
		"YV2": "trucks",
		"YV3": "buses",
		"YV4": "multipurpose",
		"4V1": "trucks",
		"4V4": "trucks",
		"4V5": "trucks",
		"MHA": "PT. Central Sole Agency",
	}
	// 车系 pos 4
	CarSeries = map[rune]string{
		'A': "[1981-1998]240/[1999-2010]2006- S80/XC90 (2015+)",
		'B': "[1981-1998]260/[1999-2010]2008- V70, XC70",
		'C': "[1999-2010]XC90 (2002-2014)",
		'D': "[1999-2010]XC60",
		'E': "[1981-1998]480",
		'F': "[1981-1998]740/[1999-2010]New S60/V60",
		'G': "[1981-1998]760/[1999-2010]V60 Plug-In Hybrid",
		'H': "[1981-1998]780",
		'J': "[1981-1998]940/[1999-2010]V70",
		'K': "[1981-1998]440, 960",
		'L': "[1981-1998]460, 850, S70/[1999-2010]1998 V70",
		'M': "[1999-2010]Volvo S40, V50, C30, C70/2013- V40",
		'R': "[1999-2010]Volvo S60",
		'S': "[1999-2010]2000- 2007 Volvo V70, XC70",
		'T': "[1999-2010]1999- 2006 Volvo S80",
		'V': "[1999-2010]Volvo V40",
	}
	// 安全系统等 pos 5
	RestraintSys = map[rune]string{
		'C': "All-New C70",
		'H': "S40 AWD, S60 AWD, S80 AWD",
		'J': "V50 AWD, V70 AWD",
		'K': "C30 FWD",
		'L': "XC60 2WD",
		'M': "XC90 5-Seater AWD/[2013]V40 Cross Country AWD",
		'N': "XC90 5-Seater FWD",
		'R': "XC90 5-Seater AWD",
		'S': "S40 FWD, S60 FWD, S80 FWD",
		'V': "V40 FWD",
		'W': "V50 FWD, V60 FWD, V60 Plug-In Hybrid (AWD), V70 FWD, V70 AWD",
		'Y': "XC90 7-Seater FWD",
		'Z': "XC60 AWD, XC70 AWD, XC90 7-Seater AWD",
	}

	// 发动机 pos 6-7
	Engine = map[string]string{
		"04": "C30/S40/V50/V70 2.0 Flexifuel FWD",
		"17": "B4204S2 V40 2.0l FWD",
		"18": "B4194T V40 1.9l Turbo FWD",
		"20": "B4164S3 C30 1.6 FWD",
		"21": "B4184S11 V50/S40 1.8",
		"30": "D5244T18 XC90 2.4 AWD",
		"38": "B5244S4 S40/V50 2.4i FWD",
		"39": "B5244S7 S40/V50 2.4i FWD",
		"40": "B4204T11 S60/V60 2.0 T5 FWD",
		"41": "B5202S 850/V70 2.0i FWD",
		"43": "B5204T3 850/V70/S80 2.0 T5 FWD",
		"47": "B5204T 850/V70 2.0 T5 FWD",
		"51": "B5252S 850/V70 2.5i FWD or D5204T6 V40/V40CC D3/D4",
		"52": "B5254T4 S60/V70 R AWD",
		"53": "B5234T3 S60/V70 T5 FWD",
		"54": "B5244T5 S60 T5 FWD",
		"55": "B5254FS 850/V70/S70 FWD",
		"56": "B5254T S70/V70 GLT FWD - 1999",
		//"56": "B5244T 2000 -",
		"57": "B5234T 850/V70 2.3T FWD Turbo",
		"58": "B5234T5 1995-1997 850 T-5R/R ; B5244T3 S60/V70/S80/XC70 2000 -",
		"59": "B5254T2 S80/S60/XC90 2.5T FWD/AWD, V70 2.5T FWD, XC70 AWD",
		"61": "B5244S S60/V70 2.4 FWD",
		"64": "B5244S6 S60/V70 2.4 FWD",
		"65": "B5244S2 S80/V70 2.4 FWD",
		"66": "B5244S5 S40/V50 2.4 FWD",
		"67": "B5244S4 C30/C70 T5 FWD",
		"68": "B5254T3 S40/V50 T5 FWD/AWD",
		"69": "D5244T5 S80/V70 2.4D FWD",
		"70": "D4192T3 S40/V40",
		//"70": "D5244T10 XC60 AWD D5(205)",
		"71": "D5244T4 V70 AWD D5(185), XC90 AWD D5(185)",
		"72": "D5252T S70/S80 2.5TDi FWD",
		"73": "D4192T2 S40/V40",
		"74": "D5244T2",
		"75": "D4204T C30 2.0D",
		"76": "D4164T 1.6D (PSA-Ford Engine)",
		"77": "D5244T8 S40 D5 AT",
		"78": "D4192T4 S40/V40",
		"79": "D5244T D5(163)",
		"82": "D5244T15 XC60 AWD D5(215)",
		"84": "D4162T S60",
		"85": "B8444S XC90/S80 V8 AWD",
		"88": "D5204T3 XC60",
		"90": "B6284T S80 2.8 T6",
		"91": "B6294T S80/XC90 2.9 T6",
		"94": "B6294S S80 2.9 FWD",
		"97": "B6299S S80 2.9 FWD",
		"98": "B6324S XC90/S80/V70 3.2 FWD/XC70 AWD",
		"99": "B6304T4 S80 3.0 T6 AWD",
		"AA": "V60 Plug-In Hybrid",
		"A9": "S60/V60 Polestar 3,0L (350)",
		"A0": "S60/V60 Polestar 2,0L (367)",
	}
	// 排放 pos 8
	EmissisonCode = map[rune]string{
		'0': "SULEV+ (Super Ultra Low Emissions Vehicle) / Engine Codes 39, 55, 64, 72",
		'2': "ULEV2 (Ultra Low Emissions Vehicle) / Engine Codes 38, 41, 51, 59, 61, 67, 68, 85, 98, 99",
		'4': "Engine Codes 71",
		'7': "LEV2 (Low Emissions Vehicle) / Engine Codes 52, 54",
		'8': "Engine Codes 70",
		'D': "L6",
		'3': "KOD: EM F1",
		'9': "KOD: EM F2",
		'B': "KOD: EM X3",
		'G': "KOD: EM X1 EXC (USA)(CDN)",
		'H': "KOD: EM X2 (USA)(CDN)",
		'J': "KOD: EM X5 EXC (USA)(CDN)",
		'Z': "KOD: EM X6 (USA)(CDN)",
		'W': "KOD: EM Z4",
	}

	// 校验位/变速箱 pos 9
	// 年份 pos 10
	// 装配厂 pos 11
	AssemblyPlant = map[rune]string{
		'0': "[Sweden] Kalmar Plant",
		'1': "[Sweden] Torslanda Plant VCT 21(Volvo Torslandaverken) (Gothenburg)",
		'2': "[Belgium] Ghent Plant VCG 22",
		'3': "[Canada] Halifax Plant",
		'4': "[Italy] - Bertone models 240",
		'5': "[Malaysia]",
		'6': "[Australia]",
		'7': "[Indonesia]",
		'A': "[Sweden] Uddevalla Plant (Volvo Cars/TWR (Tom Walkinshaw Racing))",
		'B': "[Italy] - Bertone Chongq 31",
		'D': "[Italy] - Bertone models 780",
		'E': "[Singapore]",
		'F': "[The Netherlands] Born Plant (NEDCAR)",
		'J': "[Sweden] Uddevalla Plant VCU 38 (Volvo Cars/ Pininfarina Sverige AB)",
		'M': "PVÖ 53",
	}
	// 序列号 pos 12-17
)

func GetWMIRune(vin string) vinrune.WMIRune {
	ret := vinrune.WMIRune{}
	ret.VehicleTypeStr = vin[:3]
	return ret
}

func GetVDSRune(vin string) vinrune.VDSRune {
	ret := vinrune.VDSRune{}
	vinStr := vin[3:8]
	fmt.Sscanf(vinStr, "%c%c%2s%c", &ret.CarSeriesRune, &ret.RestraintSysRune,
		&ret.EngineStr, &ret.EmissisonRune)
	return ret
}

func GetVISRune(vin string) vinrune.VISRune {
	ret := vinrune.VISRune{}
	vinStr := vin[9:]
	fmt.Sscanf(vinStr, "%c%c%6s", &ret.YearRune, &ret.AssemblyRune, &ret.SequenceNO)
	return ret
}
