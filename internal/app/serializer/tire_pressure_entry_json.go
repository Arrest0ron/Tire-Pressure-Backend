package serializer

type CoatingCoefficientUpdateJSON struct {
	CoatingCoefficient float64 `json:"coating_coefficient"`
}

type TirePressureEntryUpdateJSON struct {
	CoatingCoefficient float64 `json:"coating_coefficient"`
	Pressure           float64 `json:"pressure"`
}
