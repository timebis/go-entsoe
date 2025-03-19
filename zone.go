package entsoe

import "fmt"

type ZoneId string

const (
	// Sweden
	ZoneIdNorthSweden        ZoneId = "SE-SE1"
	ZoneIdNorthCentralSweden ZoneId = "SE-SE2"
	ZoneIdSouthCentralSweden ZoneId = "SE-SE3"
	ZoneIdSouthSweden        ZoneId = "SE-SE4"

	// Norway
	ZoneIdSoutheastNorway ZoneId = "NO-NO1"
	ZoneIdSouthwestNorway ZoneId = "NO-NO2"
	ZoneIdMiddleNorway    ZoneId = "NO-NO3"
	ZoneIdNorthNorway     ZoneId = "NO-NO4"
	ZoneIdWestNorway      ZoneId = "NO-NO5"

	// Denmark
	ZoneIdWestDenmark ZoneId = "DK-DK1"
	ZoneIdEastDenmark ZoneId = "DK-DK2"
	ZoneIdBornholm    ZoneId = "DK-BHM"

	// Finland
	ZoneIdFinland ZoneId = "FI"

	// Great Britain and Ireland
	ZoneIdGreatBritain    ZoneId = "GB"
	ZoneIdNorthernIreland ZoneId = "GB-NIR"
	ZoneIdIreland         ZoneId = "IE"

	// Western Europe
	ZoneIdFrance      ZoneId = "FR"
	ZoneIdGermany     ZoneId = "DE"
	ZoneIdNetherlands ZoneId = "NL"
	ZoneIdBelgium     ZoneId = "BE"
	ZoneIdLuxembourg  ZoneId = "LU"
	ZoneIdAustria     ZoneId = "AT"
	ZoneIdSwitzerland ZoneId = "CH"

	// Iberian Peninsula
	ZoneIdSpain    ZoneId = "ES"
	ZoneIdPortugal ZoneId = "PT"

	// Italy
	ZoneIdNorthItaly        ZoneId = "IT-NO"
	ZoneIdCentralNorthItaly ZoneId = "IT-CNO"
	ZoneIdCentralSouthItaly ZoneId = "IT-CSO"
	ZoneIdSouthItaly        ZoneId = "IT-SO"
	ZoneIdSicily            ZoneId = "IT-SIC"
	ZoneIdSardinia          ZoneId = "IT-SAR"

	// Eastern Europe
	ZoneIdPoland            ZoneId = "PL"
	ZoneIdCzechia           ZoneId = "CZ"
	ZoneIdSlovakia          ZoneId = "SK"
	ZoneIdHungary           ZoneId = "HU"
	ZoneIdSlovenia          ZoneId = "SI"
	ZoneIdCroatia           ZoneId = "HR"
	ZoneIdRomania           ZoneId = "RO"
	ZoneIdBulgaria          ZoneId = "BG"
	ZoneIdSerbia            ZoneId = "RS"
	ZoneIdBosniaHerzegovina ZoneId = "BA"

	// Baltic States
	ZoneIdEstonia   ZoneId = "EE"
	ZoneIdLatvia    ZoneId = "LV"
	ZoneIdLithuania ZoneId = "LT"

	// Nordic and Iceland
	ZoneIdIceland ZoneId = "IS"

	// Southeast Europe
	ZoneIdGreece ZoneId = "GR"
	ZoneIdCyprus ZoneId = "CY"
	ZoneIdTurkey ZoneId = "TR"

	// Asia
	ZoneIdJapanTokyo    ZoneId = "JP-TK"
	ZoneIdJapanTohoku   ZoneId = "JP-TH"
	ZoneIdJapanOkinawa  ZoneId = "JP-ON"
	ZoneIdJapanKyushu   ZoneId = "JP-KY"
	ZoneIdJapanKansai   ZoneId = "JP-KN"
	ZoneIdJapanHokuriku ZoneId = "JP-HR"
	ZoneIdJapanHokkaido ZoneId = "JP-HKD"
	ZoneIdJapanChugoku  ZoneId = "JP-CG"
	ZoneIdJapanChubu    ZoneId = "JP-CB"
	ZoneIdSouthKorea    ZoneId = "KR"
	ZoneIdTaiwan        ZoneId = "TW"
	ZoneIdHongKong      ZoneId = "HK"
	ZoneIdSingapore     ZoneId = "SG"
	ZoneIdIsrael        ZoneId = "IL"

	// India
	ZoneIdMainlandIndia     ZoneId = "IN"
	ZoneIdNorthernIndia     ZoneId = "IN-NO"
	ZoneIdWesternIndia      ZoneId = "IN-WE"
	ZoneIdSouthernIndia     ZoneId = "IN-SO"
	ZoneIdEasternIndia      ZoneId = "IN-EA"
	ZoneIdNorthEasternIndia ZoneId = "IN-NE"

	// Southeast Asia
	ZoneIdIndonesia           ZoneId = "ID"
	ZoneIdMalaysiaPeninsula   ZoneId = "MY-WM"
	ZoneIdMalaysiaBorneo      ZoneId = "MY-EM"
	ZoneIdPhilippinesLuzon    ZoneId = "PH-LU"
	ZoneIdPhilippinesVisayas  ZoneId = "PH-VI"
	ZoneIdPhilippinesMindanao ZoneId = "PH-MI"

	// Latin America
	ZoneIdUruguay         ZoneId = "UY"
	ZoneIdPeru            ZoneId = "PE"
	ZoneIdPanama          ZoneId = "PA"
	ZoneIdNicaragua       ZoneId = "NI"
	ZoneIdCostaRica       ZoneId = "CR"
	ZoneIdChileSEN        ZoneId = "CL-SEN"
	ZoneIdBrazilNorth     ZoneId = "BR-N"
	ZoneIdBrazilNorthEast ZoneId = "BR-NE"
	ZoneIdBrazilCentral   ZoneId = "BR-CS"
	ZoneIdBrazilSouth     ZoneId = "BR-S"

	// Australia
	ZoneIdAustraliaQueensland        ZoneId = "AU-QLD"
	ZoneIdAustraliaNSW               ZoneId = "AU-NSW"
	ZoneIdAustraliaVictoria          ZoneId = "AU-VIC"
	ZoneIdAustraliaSouthAustralia    ZoneId = "AU-SA"
	ZoneIdAustraliaWesternAustralia  ZoneId = "AU-WA"
	ZoneIdAustraliaTasmania          ZoneId = "AU-TAS"
	ZoneIdAustraliaNorthernTerritory ZoneId = "AU-NT"

	// North America
	ZoneIdCanadaOntario ZoneId = "CA-ON"
	ZoneIdCanadaQuebec  ZoneId = "CA-QC"

	// Africa
	ZoneIdSouthAfrica ZoneId = "ZA"
	ZoneIdKenya       ZoneId = "KE"

	// Middle East
	ZoneIdQatar ZoneId = "QA"
)

var ZoneIdToDomainType = map[ZoneId]DomainType{
	// Sweden
	ZoneIdNorthSweden:        DomainSE1,
	ZoneIdNorthCentralSweden: DomainSE2,
	ZoneIdSouthCentralSweden: DomainSE3,
	ZoneIdSouthSweden:        DomainSE4,

	// Norway
	ZoneIdSoutheastNorway: DomainNO1,
	ZoneIdSouthwestNorway: DomainNO2,
	ZoneIdMiddleNorway:    DomainNO3,
	ZoneIdNorthNorway:     DomainNO4,
	ZoneIdWestNorway:      DomainNO5,

	// Denmark
	ZoneIdWestDenmark: DomainDK1,
	ZoneIdEastDenmark: DomainDK2,
	// ZoneIdBornholm: DomainBHM, // Domain not found in types.go

	// Finland
	ZoneIdFinland: DomainFI,

	// Great Britain and Ireland
	ZoneIdGreatBritain:    DomainGB,
	ZoneIdNorthernIreland: DomainNIR,
	ZoneIdIreland:         DomainIE,

	// Western Europe
	ZoneIdFrance:      DomainFR,
	ZoneIdGermany:     DomainDE,
	ZoneIdNetherlands: DomainNL,
	ZoneIdBelgium:     DomainBE,
	ZoneIdLuxembourg:  DomainLU,
	ZoneIdAustria:     DomainAT,
	ZoneIdSwitzerland: DomainCH,

	// Iberian Peninsula
	ZoneIdSpain:    DomainES,
	ZoneIdPortugal: DomainPT,

	// Italy
	ZoneIdNorthItaly:        DomainITNorth,
	ZoneIdCentralNorthItaly: DomainITZCentreNorth,
	ZoneIdCentralSouthItaly: DomainITCentreSouth,
	ZoneIdSouthItaly:        DomainITZSouth,
	ZoneIdSicily:            DomainITSicily,
	ZoneIdSardinia:          DomainITZSardinia,

	// Eastern Europe
	ZoneIdPoland:            DomainPL,
	ZoneIdCzechia:           DomainCZ,
	ZoneIdSlovakia:          DomainSK,
	ZoneIdHungary:           DomainHU,
	ZoneIdSlovenia:          DomainSI,
	ZoneIdCroatia:           DomainHR,
	ZoneIdRomania:           DomainRO,
	ZoneIdBulgaria:          DomainBG,
	ZoneIdSerbia:            DomainRS,
	ZoneIdBosniaHerzegovina: DomainBA,

	// Baltic States
	ZoneIdEstonia:   DomainEE,
	ZoneIdLatvia:    DomainLV,
	ZoneIdLithuania: DomainLT,

	// Nordic and Iceland
	ZoneIdIceland: DomainIS,

	// Southeast Europe
	ZoneIdGreece: DomainGR,
	ZoneIdCyprus: DomainCY,
	ZoneIdTurkey: DomainTR,
}

func GetDomainType(zone ZoneId) (DomainType, error) {
	domain, ok := ZoneIdToDomainType[zone]
	if !ok {
		return "", fmt.Errorf("unsupported area %s", zone)
	}

	return domain, nil
}
