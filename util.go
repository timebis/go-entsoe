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
		PsrTypeFossilBrownCoalLignite:     "FossilBrownCoalLignite",
		PsrTypeFossilCoalDerivedGas:       "FossilCoalDerivedGas",
		PsrTypeFossilGas:                  "FossilGas",
		PsrTypeFossilHardCoal:             "FossilHardCoal",
		PsrTypeFossilOil:                  "FossilOil",
		PsrTypeFossilOilShale:             "FossilOilShale",
		PsrTypeFossilPeat:                 "FossilPeat",
		PsrTypeGeothermal:                 "Geothermal",
		PsrTypeHydroPumpedStorage:         "HydroPumpedStorage",
		PsrTypeHydroRunOfRiverAndPoundage: "HydroRunOfRiverAndPoundage",
		PsrTypeHydroWaterReservoir:        "HydroWaterReservoir",
		PsrTypeMarine:                     "Marine",
		PsrTypeNuclear:                    "Nuclear",
		PsrTypeOtherRenewable:             "OtherRenewable",
		PsrTypeSolar:                      "Solar",
		PsrTypeWaste:                      "Waste",
		PsrTypeWindOffshore:               "WindOffshore",
		PsrTypeWindOnshore:                "WindOnshore",
		PsrTypeOther:                      "Other",
		PsrTypeACLink:                     "ACLink",
		PsrTypeDCLink:                     "DCLink",
		PsrTypeSubstation:                 "Substation",
		PsrTypeTransformer:                "Transformer",
	}
	return psrTypeMap[psrType]
}
