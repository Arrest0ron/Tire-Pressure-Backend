package serializer

// TirePressureEntryUpdateJSON — editable fields in a tire-pressure-tire link.
type TirePressureEntryUpdateJSON struct {
	CoatingCoefficient float64 `json:"coating_coefficient"`
	Pressure           float64 `json:"pressure"`
}
