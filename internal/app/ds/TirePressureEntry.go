package ds

type TirePressureEntry struct {
	ID             uint         `gorm:"primaryKey"`
	TirePressureID uint         `gorm:"not null"`
	TirePressure   TirePressure `gorm:"foreignKey:TirePressureID"`
	TireID         uint         `gorm:"not null"`
	// ✅ ДОБАВЬТЕ references:TireID - это критично!
	Tire               Tire `gorm:"foreignKey:TireID;references:TireID"`
	CoatingCoefficient float64
	Pressure           float64
}

func (TirePressureEntry) TableName() string {
	return "tire_pressure_entries"
}
