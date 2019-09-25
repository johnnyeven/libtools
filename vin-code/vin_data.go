package vincode

// WMI 世界制造厂识别代码
type WMIData struct {
	// 洲
	Continent string
	// 制造国
	Country string
	// 制造商
	Manufacturer string
	// 车用途类型
	VehicleType string
}

// VDS 车辆说明部分
type VDSData struct {
	// 约束系统
	RestraintSystem string
	// 系列
	CarSeries string
	// 几门类型
	DoorType string
	// 发动机
	Engine string
	// 排放
	Emissison string
}

// VIS 车辆指示部分
type VISData struct {
	// 生产年份
	ModelYear string
	// 装配厂
	AssemblyPlant string
	// 生产序列号
	SequenceNO string
}

// VIN 解析后数据
type VINData struct {
	// WMI部分
	WMISection WMIData
	// VDS部分
	VDSSection VDSData
	// VIS部分
	VISSection VISData
}

type VINInterface interface {
	// 解析WMI
	ParseWMI() (WMIData, error)
	// 解析VDS
	ParseVDS() (VDSData, error)
	// 解析VIS
	ParseVIS() (VISData, error)
}

// 解析ALL
func ParseAll(it VINInterface) (VINData, error) {
	vinData := VINData{}
	WMISection, err := it.ParseWMI()
	if err != nil {
		return vinData, err
	}
	vinData.WMISection = WMISection
	VISSection, err := it.ParseVIS()
	if err != nil {
		return vinData, err
	}
	vinData.VISSection = VISSection
	VDSSection, err := it.ParseVDS()
	if err != nil {
		return vinData, err
	}
	vinData.VDSSection = VDSSection

	return vinData, nil
}
