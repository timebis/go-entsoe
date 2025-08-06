package entsoe

import (
	"fmt"
	"strings"
)

// https://transparency.entsoe.eu/content/static_content/Static%20content/web%20api/Guide.html#_areas

type Area string

const (
	// Central Western Europe
	Austria     Area = "AT"
	Belgium     Area = "BE"
	France      Area = "FR"
	Germany     Area = "DE"
	Netherlands Area = "NL"
	Poland      Area = "PL"

	// Nordic
	Denmark1 Area = "DK1"
	Denmark2 Area = "DK2"
	Finland  Area = "FI"
	Norway1  Area = "NO1"
	Norway2  Area = "NO2"
	Norway3  Area = "NO3"
	Norway4  Area = "NO4"
	Norway5  Area = "NO5"
	Sweden1  Area = "SE1"
	Sweden2  Area = "SE2"
	Sweden3  Area = "SE3"
	Sweden4  Area = "SE4"

	// Baltic
	Estonia   Area = "EE"
	Lithuania Area = "LT"
	Latvia    Area = "LV"
)

var domains = map[Area]DomainType{
	// Central Western Europe
	Austria:     DomainAT,
	Belgium:     DomainBE,
	France:      DomainFR,
	Germany:     DomainDELU,
	Netherlands: DomainNL,
	Poland:      DomainPL,

	// Nordic
	Denmark1: DomainDK1,
	Denmark2: DomainDK2,
	Finland:  DomainFI,
	Norway1:  DomainNO1,
	Norway2:  DomainNO2,
	Norway3:  DomainNO3,
	Norway4:  DomainNO4,
	Norway5:  DomainNO5,
	Sweden1:  DomainSE1,
	Sweden2:  DomainSE2,
	Sweden3:  DomainSE3,
	Sweden4:  DomainSE4,

	// Baltic
	Estonia:   DomainEE,
	Lithuania: DomainLT,
	Latvia:    DomainLV,
}

var contryNames = map[Area]string{
	// Central Western Europe
	Austria:     "Austria",
	Belgium:     "Belgium",
	France:      "France",
	Germany:     "Germany",
	Netherlands: "Netherlands",
	Poland:      "Poland",

	// Nordic
	Denmark1: "Denmark",
	Denmark2: "Denmark",
	Finland:  "Finland",
	Norway1:  "Norway",
	Norway2:  "Norway",
	Norway3:  "Norway",
	Norway4:  "Norway",
	Norway5:  "Norway",
	Sweden1:  "Sweden",
	Sweden2:  "Sweden",
	Sweden3:  "Sweden",
	Sweden4:  "Sweden",

	// Baltic
	Estonia:   "Estonia",
	Lithuania: "Lithuania",
	Latvia:    "Latvia",
}

func domain(area string) (DomainType, error) {
	zone := Area(strings.ToUpper(area))

	domain, ok := domains[zone]
	if !ok {
		return "", fmt.Errorf("unsupported area %s", area)
	}

	return domain, nil
}

func (a Area) FullName() (string, error) {
	fullName, ok := contryNames[a]
	if !ok {
		return "", fmt.Errorf("unsupported area %s", a)
	}

	return fullName, nil
}
