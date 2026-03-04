package utils

// CalculatePressure calculates the tire pressure based on coating coefficient
// This is a simplified example - in a real system this would use actual tire pressure formulas
func CalculatePressure(coatingCoeff float64) float64 {
	// Base pressure calculation formula
	// This could be more complex based on tire specifications, vehicle weight, etc.
	basePressure := 2.2 // Base pressure in bar
	return basePressure * coatingCoeff
}
