package ds

type TirePressureEntry  struct {
	ID                 uint   `gorm:"primaryKey"`
	TirePressureID uint   `gorm:"not null;uniqueIndex:idx_tire_pressure_entry"`
	TireID     uint   `gorm:"not null;uniqueIndex:idx_tire_pressure_entry"`
	CoatingCoeff     float64       `gorm:"default:null"`


	TirePressure TirePressure `gorm:"foreignKey:TirePressureID"`
	Tire    Tire     `gorm:"foreignKey:TireID"`
}

// TableName — таблица связи заявка–конструкция в БД.
func (TirePressureEntry) TableName() string {
	return "tire_pressure_entries"
}
