package misc

// 洲名
const (
	AfricaStr       = "Africa"
	AsiaStr         = "Asia"
	EuropeStr       = "Europe"
	NorthAmericaStr = "North America"
	OceaniaStr      = "Oceania"
	SouthAmericaStr = "South America"

	ContinentUnknownStr = "Unknown continent"
)

// 国名
const (
	SouthAfricaNm    = "South Africa"
	IvoryCoastNm     = "Ivory Coast"
	AfricaNm         = "Africa"
	AngolaNm         = "Angola"
	KenyaNm          = "Kenya"
	TanzaniaNm       = "Tanzania"
	BeninNm          = "Benin"
	MadagascarNm     = "Madagascar"
	TunisiaNm        = "Tunisia"
	EgyptNm          = "Egypt"
	MoroccoNm        = "Morocco"
	ZambiaNm         = "Zambia"
	GhanaNm          = "Ghana"
	MozambiqueNm     = "Mozambique"
	NigeriaNm        = "Nigeria"
	JapanNm          = "Japan"
	SriLankaNm       = "SriLanka"
	IsraelNm         = "Israel"
	SKoreaNm         = "SKorea"
	KazakhstanNm     = "Kazakhstan"
	ChinaNm          = "China"
	IndiaNm          = "India"
	IndonesiaNm      = "Indonesia"
	ThailandNm       = "Thailand"
	IranNm           = "Iran"
	PakistanNm       = "Pakistan"
	TurkeyNm         = "Turkey"
	PhilippinesNm    = "Philippines"
	SingaporeNm      = "Singapore"
	MalaysiaNm       = "Malaysia"
	AsiaNm           = "Asia"
	UAENm            = "UAE"
	TaiwanNm         = "Taiwan"
	VietnamNm        = "Vietnam"
	SaudiArabiaNm    = "Saudi Arabia"
	UnitedKingdomNm  = "United Kingdom"
	EGermanyNm       = "EGermany"
	PolandNm         = "Poland"
	LatviaNm         = "Latvia"
	SwitzerlandNm    = "Switzerland"
	CzechRepNm       = "Czech Rep"
	HungaryNm        = "Hungary"
	PortugalNm       = "Portugal"
	DenmarkNm        = "Denmark"
	IrelandNm        = "Ireland"
	RomaniaNm        = "Romania"
	SlovakNm         = "Slovak"
	AustriaNm        = "Austria"
	FranceNm         = "France"
	SpainNm          = "Spain"
	SerbiaNm         = "Serbia"
	CroatiaNm        = "Croatia"
	EstoniaNm        = "Estonia"
	GermanyNm        = "Germany"
	BulgariaNm       = "Bulgaria"
	GreeceNm         = "Greece"
	NetherlandsNm    = "Netherlands"
	USSRNm           = "USSR"
	LuxembourgNm     = "Luxembourg"
	RussiaNm         = "Russia"
	BelgiumNm        = "Belgium"
	FinlandNm        = "Finland"
	MaltaNm          = "Malta"
	SwedenNm         = "Sweden"
	NorwayNm         = "Norway"
	BelarusNm        = "Belarus"
	UkraineNm        = "Ukraine"
	ItalyNm          = "Italy"
	SloveniaNm       = "Slovenia"
	LithuaniaNm      = "Lithuania"
	EuropeNm         = "Europe"
	CanadaNm         = "Canada"
	MexicoNm         = "Mexico"
	UnitedStatesNm   = "UnitedStates"
	AustraliaNm      = "Australia"
	NewZealandNm     = "NewZealand"
	OceaniaNm        = "Oceania"
	ArgentinaNm      = "Argentina"
	ChileNm          = "Chile"
	EcuadorNm        = "Ecuador"
	PeruNm           = "Peru"
	VenezuelaNm      = "Venezuela"
	SouthAmericaNm   = "SouthAmerica"
	ColombiaNm       = "Colombia"
	ParaguayNm       = "Paraguay"
	UruguayNm        = "Uruguay"
	TrinidadTobagoNm = "Trinidad & Tobago"
	BrazilNm         = "Brazil"
	SANm             = "SA"

	CountryUnknownNm = "Unknown Country"
)

// 获取洲属
func GetVINContinent(vin string) string {
	r := vin[0]
	switch {
	case r >= 'A' && r <= 'H':
		return AfricaStr
	case r >= 'J' && r <= 'R':
		return AsiaStr
	case r >= 'S' && r <= 'Z':
		return EuropeStr
	case r >= '1' && r <= '5':
		return NorthAmericaStr
	case r >= '6' && r <= '7':
		return OceaniaStr
	case r == '8' || r == '9' || r == '0':
		return SouthAmericaStr
	default:
		return ContinentUnknownStr
	}
}

func GetAfricaCountry(vin string) string {
	nm1 := vin[0]
	nm2 := vin[1]
	switch {
	case nm1 == 'A' && (nm2 >= 'A' && nm2 <= 'H'):
		return SouthAfricaNm
	case nm1 == 'A' && (nm2 >= 'J' && nm2 <= 'N'):
		return IvoryCoastNm
	case nm1 == 'B' && (nm2 >= 'A' && nm2 <= 'E'):
		return AngolaNm
	case nm1 == 'C' && (nm2 >= 'A' && nm2 <= 'E'):
		return BeninNm
	case nm1 == 'D' && (nm2 >= 'A' && nm2 <= 'E'):
		return EgyptNm
	case nm1 == 'E' && (nm2 >= 'A' && nm2 <= 'E'):
		return GhanaNm
	case nm1 == 'F' && (nm2 >= 'A' && nm2 <= 'E'):
		return IndiaNm
	case nm1 == 'B' && (nm2 >= 'F' && nm2 <= 'K'):
		return KenyaNm
	case nm1 == 'C' && (nm2 >= 'F' && nm2 <= 'K'):
		return MadagascarNm
	case nm1 == 'D' && (nm2 >= 'F' && nm2 <= 'K'):
		return MoroccoNm
	case nm1 == 'E' && (nm2 >= 'F' && nm2 <= 'K'):
		return MozambiqueNm
	case nm1 == 'F' && (nm2 >= 'F' && nm2 <= 'K'):
		return NigeriaNm
	case nm1 == 'B' && (nm2 >= 'L' && nm2 <= 'R'):
		return TanzaniaNm
	case nm1 == 'C' && (nm2 >= 'L' && nm2 <= 'R'):
		return TunisiaNm
	case nm1 == 'D' && (nm2 >= 'L' && nm2 <= 'R'):
		return ZambiaNm
	default:
		return AfricaNm
	}
}

func GetAsiaCountry(vin string) string {
	nm1 := vin[0]
	nm2 := vin[1]
	switch {
	case nm1 == 'J':
		return JapanNm
	case nm1 == 'L':
		return ChinaNm
	case nm1 == 'K' && (nm2 >= 'A' && nm2 <= 'E'):
		return SriLankaNm
	case nm1 == 'M' && (nm2 >= 'A' && nm2 <= 'E'):
		return IndiaNm
	case nm1 == 'N' && (nm2 >= 'A' && nm2 <= 'E'):
		return IranNm
	case nm1 == 'P' && (nm2 >= 'A' && nm2 <= 'E'):
		return PhilippinesNm
	case nm1 == 'R' && (nm2 >= 'A' && nm2 <= 'E'):
		return UAENm
	case nm1 == 'K' && (nm2 >= 'F' && nm2 <= 'K'):
		return IsraelNm
	case nm1 == 'M' && (nm2 >= 'F' && nm2 <= 'K'):
		return IndonesiaNm
	case nm1 == 'N' && (nm2 >= 'F' && nm2 <= 'K'):
		return PakistanNm
	case nm1 == 'P' && (nm2 >= 'F' && nm2 <= 'K'):
		return SingaporeNm
	case nm1 == 'R' && (nm2 >= 'F' && nm2 <= 'K'):
		return TaiwanNm
	case nm1 == 'K' && (nm2 >= 'L' && nm2 <= 'R'):
		return SKoreaNm
	case nm1 == 'M' && (nm2 >= 'L' && nm2 <= 'R'):
		return ThailandNm
	case nm1 == 'N' && (nm2 >= 'L' && nm2 <= 'R'):
		return TurkeyNm
	case nm1 == 'P' && (nm2 >= 'L' && nm2 <= 'R'):
		return MalaysiaNm
	case nm1 == 'R' && (nm2 >= 'L' && nm2 <= 'R'):
		return VietnamNm
	case nm1 == 'R' && ((nm2 >= 'S' && nm2 <= 'Z') || (nm2 >= '0' && nm2 <= '9')):
		return SaudiArabiaNm
	default:
		return AsiaNm
	}
}
func GetEuropeCountry(vin string) string {
	nm1 := vin[0]
	nm2 := vin[1]
	switch {
	case nm1 == 'S' && (nm2 >= 'A' && nm2 <= 'M'):
		return UnitedKingdomNm
	case nm1 == 'S' && (nm2 >= 'N' && nm2 <= 'T'):
		return EGermanyNm
	case nm1 == 'S' && (nm2 >= 'U' && nm2 <= 'Z'):
		return PolandNm
	case nm1 == 'S' && (nm2 >= '1' && nm2 <= '4'):
		return LatviaNm
	case nm1 == 'T' && (nm2 >= 'A' && nm2 <= 'H'):
		return SwitzerlandNm
	case nm1 == 'T' && (nm2 >= 'J' && nm2 <= 'P'):
		return CzechRepNm
	case nm1 == 'T' && (nm2 >= 'R' && nm2 <= 'V'):
		return HungaryNm
	case nm1 == 'T' && ((nm2 >= 'W' && nm2 <= 'Z') || nm2 == '1'):
		return PortugalNm
	case nm1 == 'U' && (nm2 >= 'H' && nm2 <= 'M'):
		return DenmarkNm
	case nm1 == 'U' && (nm2 >= 'N' && nm2 <= 'T'):
		return IrelandNm
	case nm1 == 'U' && (nm2 >= 'U' && nm2 <= 'Z'):
		return RomaniaNm
	case nm1 == 'U' && (nm2 >= '5' && nm2 <= '7'):
		return SlovakNm
	case nm1 == 'V' && (nm2 >= 'A' && nm2 <= 'E'):
		return SlovakNm
	case nm1 == 'V' && (nm2 >= 'F' && nm2 <= 'R'):
		return AustriaNm
	case nm1 == 'V' && (nm2 >= 'F' && nm2 <= 'R'):
		return FranceNm
	case nm1 == 'V' && (nm2 >= 'S' && nm2 <= 'W'):
		return SpainNm
	case nm1 == 'V' && ((nm2 >= 'X' && nm2 <= 'Z') || (nm2 >= '1' || nm2 <= '2')):
		return SerbiaNm
	case nm1 == 'V' && (nm2 >= '3' && nm2 <= '5'):
		return CroatiaNm
	case nm1 == 'V' && ((nm2 >= '6' && nm2 <= '9') || (nm2 == '0')):
		return EstoniaNm
	case nm1 == 'W':
		return GermanyNm
	case nm1 == 'X' && (nm2 >= 'A' && nm2 <= 'E'):
		return BulgariaNm
	case nm1 == 'X' && (nm2 >= 'F' && nm2 <= 'K'):
		return GreeceNm
	case nm1 == 'X' && (nm2 >= 'L' && nm2 <= 'R'):
		return NetherlandsNm
	case nm1 == 'X' && (nm2 >= 'S' && nm2 <= 'W'):
		return USSRNm
	case nm1 == 'X' && ((nm2 >= 'X' && nm2 <= 'Z') || (nm2 == '1' || nm2 == '2')):
		return LuxembourgNm
	case nm1 == 'X' && ((nm2 >= '3' && nm2 <= '9') || (nm2 == '0')):
		return RussiaNm
	case nm1 == 'Y' && (nm2 >= 'A' && nm2 <= 'E'):
		return BelgiumNm
	case nm1 == 'Y' && (nm2 >= 'F' && nm2 <= 'K'):
		return FinlandNm
	case nm1 == 'Y' && (nm2 >= 'L' && nm2 <= 'R'):
		return MaltaNm
	case nm1 == 'Y' && (nm2 >= 'S' && nm2 <= 'W'):
		return SwedenNm
	case nm1 == 'Y' && ((nm2 >= 'X' && nm2 <= 'Z') || (nm2 == '1' || nm2 == '2')):
		return NorwayNm
	case nm1 == 'Y' && (nm2 >= '3' && nm2 <= '5'):
		return BelarusNm
	case nm1 == 'Y' && ((nm2 >= '6' && nm2 <= '9') || (nm2 == '0')):
		return UkraineNm
	case nm1 == 'Z' && (nm2 >= 'A' && nm2 <= 'R'):
		return ItalyNm
	case nm1 == 'Z' && ((nm2 >= 'X' && nm2 <= 'Z') || (nm2 == '1' || nm2 == '2')):
		return SloveniaNm
	case nm1 == 'Z' && (nm2 >= '3' && nm2 <= '5'):
		return LithuaniaNm
	default:
		return EuropeStr
	}
}

func GetNorthAmericaCountry(vin string) string {
	nm1 := vin[0]
	switch nm1 {
	case '1', '4', '5':
		return UnitedStatesNm
	case '2':
		return CanadaNm
	case '3':
		return MexicoNm
	default:
		return ""
	}
}

func GetOceaniaCountry(vin string) string {
	nm1 := vin[0]
	nm2 := vin[1]
	switch {
	case nm1 == '6' && (nm2 >= 'A' && nm2 <= 'W'):
		return AustriaNm
	case nm1 == '7' && (nm2 >= 'A' && nm2 <= 'E'):
		return NewZealandNm
	default:
		return OceaniaNm
	}
}

func GetSouthAmericaCountry(vin string) string {
	nm1 := vin[0]
	nm2 := vin[1]
	switch {
	case nm1 == '8' && (nm2 >= 'A' && nm2 <= 'E'):
		return ArgentinaNm
	case nm1 == '8' && (nm2 >= 'F' && nm2 <= 'K'):
		return ChileNm
	case nm1 == '8' && (nm2 >= 'L' && nm2 <= 'R'):
		return EcuadorNm
	case nm1 == '8' && (nm2 >= 'S' && nm2 <= 'W'):
		return PeruNm
	case nm1 == '8' && ((nm2 >= 'X' && nm2 <= 'Z') || (nm2 == '1' || nm2 == '2')):
		return VenezuelaNm
	case nm1 == '9' && (nm2 >= 'A' && nm2 <= 'E'):
		return BrazilNm
	case nm1 == '9' && (nm2 >= 'F' && nm2 <= 'K'):
		return ColombiaNm
	case nm1 == '9' && (nm2 >= 'L' && nm2 <= 'R'):
		return ParaguayNm
	case nm1 == '9' && (nm2 >= 'S' && nm2 <= 'W'):
		return UruguayNm
	case nm1 == '9' && ((nm2 >= 'X' && nm2 <= 'Z') || (nm2 == '1' || nm2 == '2')):
		return TrinidadTobagoNm
	case nm1 == '9' && (nm2 >= '3' && nm2 <= '9'):
		return BrazilNm
	case nm1 == '9' && (nm2 == '0'):
		return SANm
	default:
		return SouthAmericaNm
	}
}

// 获取国属
func GetVINCountry(vin string) string {
	conti := GetVINContinent(vin)
	switch conti {
	case AfricaStr:
		return GetAfricaCountry(vin)
	case AsiaStr:
		return GetAsiaCountry(vin)
	case EuropeStr:
		return GetEuropeCountry(vin)
	case NorthAmericaStr:
		return GetNorthAmericaCountry(vin)
	case OceaniaStr:
		return GetOceaniaCountry(vin)
	case SouthAmericaStr:
		return GetSouthAmericaCountry(vin)
	default:
		return ""
	}
}
