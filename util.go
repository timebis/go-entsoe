package entsoe
// GetPsrTypeName returns the human-readable name of a PsrType.
// For example, GetPsrTypeName(PsrTypeNuclear) returns "Nuclear".
// Returns an empty string if the PsrType is not recognized.
func GetPsrTypeName(psrType PsrType) string {
	psrTypeMap := map[PsrType]string{
		PsrTypeMixed:                      "Mixed",
		PsrTypeGeneration:                 "Generation",
		PsrTypeLoad:                       "Load",
		PsrTypeBiomass:                    "Biomass",
		PsrTypeFossilBrownCoalLignite:     "Fossil Brown Coal Lignite",
		PsrTypeFossilCoalDerivedGas:       "Fossil Coal Derived Gas",
		PsrTypeFossilGas:                  "Fossil Gas",
		PsrTypeFossilHardCoal:             "Fossil Hard Coal",
		PsrTypeFossilOil:                  "Fossil Oil",
		PsrTypeFossilOilShale:             "Fossil Oil Shale",
		PsrTypeFossilPeat:                 "Fossil Peat",
		PsrTypeGeothermal:                 "Geothermal",
		PsrTypeHydroPumpedStorage:         "Hydro Pumped Storage",
		PsrTypeHydroRunOfRiverAndPoundage: "Hydro Run Of River And Poundage",
		PsrTypeHydroWaterReservoir:        "Hydro Water Reservoir",
		PsrTypeMarine:                     "Marine",
		PsrTypeNuclear:                    "Nuclear",
		PsrTypeOtherRenewable:             "Other Renewable",
		PsrTypeSolar:                      "Solar",
		PsrTypeWaste:                      "Waste",
		PsrTypeWindOffshore:               "Wind Offshore",
		PsrTypeWindOnshore:                "Wind Onshore",
		PsrTypeOther:                      "Other",
		PsrTypeACLink:                     "ACLink",
		PsrTypeDCLink:                     "DCLink",
		PsrTypeSubstation:                 "Substation",
		PsrTypeTransformer:                "Transformer",
	}
	return psrTypeMap[psrType]
}
