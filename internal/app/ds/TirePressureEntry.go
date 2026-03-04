package ds

type TirePressureEntry struct {
	ID               int           `gorm:"primaryKey" json:"id"`
	TirePressureID   int           `gorm:"not null;index" json:"tire_pressure_id"`
	TireID           int           `gorm:"not null;index" json:"tire_id"`
	CoatingCoeff     float64       `gorm:"not null" json:"coating_coeff"`
	Pressure         float64       `gorm:"not null" json:"pressure"`
	TirePressure     TirePressure  `gorm:"foreignKey:TirePressureID" json:"-"`
	Tire             Tire          `gorm:"foreignKey:TireID" json:"-"`

	// Composite unique key to prevent duplicate entries
	// gorm:"uniqueIndex:idx_tire_pressure_tire"`
}
